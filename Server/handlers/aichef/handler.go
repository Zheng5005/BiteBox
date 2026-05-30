package aichef

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Zheng5005/BiteBox/internal/ai"
	"github.com/Zheng5005/BiteBox/internal/ratelimit"
	"github.com/Zheng5005/BiteBox/internal/recipes"
	"github.com/Zheng5005/BiteBox/utils"
)

// AIChefHandler handles AI Chef recipe generation requests.
type AIChefHandler struct {
	CatalogRepo *recipes.CatalogRepo
	AIClient    *ai.AIClient
	RateLimiter *ratelimit.RateLimiter
	Threshold   int
}

// NewAIChefHandler creates a new AI Chef handler.
func NewAIChefHandler(catalog *recipes.CatalogRepo, aiClient *ai.AIClient, rl *ratelimit.RateLimiter, threshold int) *AIChefHandler {
	return &AIChefHandler{
		CatalogRepo: catalog,
		AIClient:    aiClient,
		RateLimiter: rl,
		Threshold:   threshold,
	}
}

// AiChefHandler handles POST /api/ai/recipes with SSE streaming.
func (h *AIChefHandler) AiChefHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.ParseToken(r, "")
	if err != nil {
		h.sendError(w, "Authentication required", "UNAUTHORIZED")
		return
	}

	var req AiChefRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", "INVALID_INPUT")
		return
	}

	if len(req.Ingredients) == 0 || len(req.Ingredients) > 20 {
		h.sendError(w, "Ingredients must be 1-20 items", "INVALID_INPUT")
		return
	}

	tags := recipes.NormalizeIngredients(req.Ingredients)

	if !h.RateLimiter.Allow(userID) {
		h.sendError(w, "Daily AI request limit reached (10/day)", "RATE_LIMITED")
		return
	}

	wantsSSE := !strings.Contains(r.Header.Get("Accept"), "application/json") ||
		r.Header.Get("Accept") == "" ||
		r.Header.Get("Accept") == "*/*"

	if wantsSSE {
		h.handleSSE(w, r, userID, req, tags)
	} else {
		h.handleJSON(w, r, userID, req, tags)
	}
}

func (h *AIChefHandler) handleSSE(w http.ResponseWriter, r *http.Request, userID string, req AiChefRequest, tags []string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.sendErrorPlain(w, "Streaming not supported", "INTERNAL_ERROR")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var searchConstraints *recipes.SearchConstraints
	if req.Constraints != nil {
		searchConstraints = &recipes.SearchConstraints{
			MaxTime:    req.Constraints.MaxTime,
			Difficulty: req.Constraints.Difficulty,
			Cuisine:    req.Constraints.Cuisine,
		}
	}

	// Step 1: Catalog search
	catalogResults, err := h.CatalogRepo.SearchByIngredients(ctx, tags, searchConstraints)
	if err != nil {
		h.sendSSEError(w, "Catalog search failed", "INTERNAL_ERROR", flusher)
		return
	}

	// Step 2: Threshold gate
	if len(catalogResults) >= h.Threshold {
		event := CatalogResultEvent{
			Source:     "catalog",
			Recipes:    toRecipeSummaries(catalogResults),
			MatchCount: len(catalogResults),
		}
		writeSSEEvent(w, "catalog_results", event, flusher)
		h.recordInteraction(ctx, userID, req.Ingredients, len(catalogResults), nil, false)
		return
	}

	// Step 3: AI generation path
	writeSSEEvent(w, "ai_generation", AiGenerationEvent{
		Source: "ai",
		Status: "starting",
	}, flusher)

	catalogHints := make([]string, len(catalogResults))
	for i, cr := range catalogResults {
		catalogHints[i] = cr.Name
	}

	aiConstraints := toAiChefConstraints(req.Constraints)
	prompt := ai.BuildPrompt(req.Ingredients, aiConstraints, catalogHints)

	// Stream AI generation
	tokenCh, errCh := h.AIClient.GenerateStream(ctx, prompt)

	var accumulated string
	done := make(chan struct{})
	var streamErr error

	go func() {
		defer close(done)
		for {
			select {
			case token, ok := <-tokenCh:
				if !ok {
					return
				}
				accumulated += token
				writeSSEEvent(w, "ai_generation", AiGenerationEvent{
					Source: "ai",
					Status: "partial",
					Token:  token,
				}, flusher)
			case err := <-errCh:
				streamErr = err
				return
			}
		}
	}()

	<-done

	if streamErr != nil {
		h.sendSSEError(w, "AI generation failed: "+streamErr.Error(), "INTERNAL_ERROR", flusher)
		return
	}

	parsedRecipe, err := ai.ParseAndValidate(accumulated)
	if err != nil {
		h.sendSSEError(w, "Invalid AI response: "+err.Error(), "INTERNAL_ERROR", flusher)
		return
	}

	// Step 4: Deduplication
	existing, dedupErr := h.CatalogRepo.FindNearDuplicate(ctx, parsedRecipe)
	if dedupErr != nil {
		log.Printf("dedup error: %v", dedupErr)
	}

	if existing != nil {
		detail, err := h.CatalogRepo.GetAiRecipeDetail(ctx, existing.ID)
		if err != nil {
			h.sendSSEError(w, "Failed to fetch existing recipe", "INTERNAL_ERROR", flusher)
			return
		}
		createdAt := h.CatalogRepo.GetRecipeCreatedAt(ctx, existing.ID)
		writeSSEEvent(w, "ai_generation", AiGenerationEvent{
			Source: "ai",
			Status: "complete",
			Recipe: toCompleteRecipe(detail, createdAt),
		}, flusher)
		h.recordInteraction(ctx, userID, req.Ingredients, 1, &existing.ID, false)
		return
	}

	// Step 5: Persist
	recipeID, err := h.CatalogRepo.SaveAiRecipe(ctx, parsedRecipe, prompt, tags)
	if err != nil {
		h.sendSSEError(w, "Failed to save recipe", "INTERNAL_ERROR", flusher)
		return
	}

	recipeIDStr := strconv.Itoa(recipeID)
	detail, err := h.CatalogRepo.GetAiRecipeDetail(ctx, recipeIDStr)
	if err != nil {
		h.sendSSEError(w, "Recipe saved but failed to fetch detail", "INTERNAL_ERROR", flusher)
		return
	}
	createdAt := h.CatalogRepo.GetRecipeCreatedAt(ctx, recipeIDStr)

	writeSSEEvent(w, "ai_generation", AiGenerationEvent{
		Source: "ai",
		Status: "complete",
		Recipe: toCompleteRecipe(detail, createdAt),
	}, flusher)

	h.recordInteraction(ctx, userID, req.Ingredients, 1, &recipeIDStr, true)
}

func (h *AIChefHandler) handleJSON(w http.ResponseWriter, r *http.Request, userID string, req AiChefRequest, tags []string) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var searchConstraints *recipes.SearchConstraints
	if req.Constraints != nil {
		searchConstraints = &recipes.SearchConstraints{
			MaxTime:    req.Constraints.MaxTime,
			Difficulty: req.Constraints.Difficulty,
			Cuisine:    req.Constraints.Cuisine,
		}
	}

	catalogResults, err := h.CatalogRepo.SearchByIngredients(ctx, tags, searchConstraints)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorEvent{Error: "Catalog search failed", Code: "INTERNAL_ERROR"})
		return
	}

	if len(catalogResults) >= h.Threshold {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CatalogResultEvent{
			Source:     "catalog",
			Recipes:    toRecipeSummaries(catalogResults),
			MatchCount: len(catalogResults),
		})
		h.recordInteraction(ctx, userID, req.Ingredients, len(catalogResults), nil, false)
		return
	}

	catalogHints := make([]string, len(catalogResults))
	for i, cr := range catalogResults {
		catalogHints[i] = cr.Name
	}

	aiConstraints := toAiChefConstraints(req.Constraints)
	prompt := ai.BuildPrompt(req.Ingredients, aiConstraints, catalogHints)

	parsedRecipe, err := h.AIClient.Generate(ctx, prompt)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			w.WriteHeader(http.StatusGatewayTimeout)
			json.NewEncoder(w).Encode(ErrorEvent{Error: "AI generation timed out", Code: "AI_TIMEOUT"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorEvent{Error: "AI generation failed: " + err.Error(), Code: "INTERNAL_ERROR"})
		return
	}

	existing, _ := h.CatalogRepo.FindNearDuplicate(ctx, parsedRecipe)
	if existing != nil {
		detail, err := h.CatalogRepo.GetAiRecipeDetail(ctx, existing.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorEvent{Error: "Failed to fetch recipe", Code: "INTERNAL_ERROR"})
			return
		}
		createdAt := h.CatalogRepo.GetRecipeCreatedAt(ctx, existing.ID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AiGenerationEvent{
			Source: "ai",
			Status: "complete",
			Recipe: toCompleteRecipe(detail, createdAt),
		})
		h.recordInteraction(ctx, userID, req.Ingredients, 1, &existing.ID, false)
		return
	}

	recipeID, err := h.CatalogRepo.SaveAiRecipe(ctx, parsedRecipe, prompt, tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorEvent{Error: "Failed to save recipe", Code: "INTERNAL_ERROR"})
		return
	}

	recipeIDStr := strconv.Itoa(recipeID)
	detail, err := h.CatalogRepo.GetAiRecipeDetail(ctx, recipeIDStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorEvent{Error: "Recipe saved but detail fetch failed", Code: "INTERNAL_ERROR"})
		return
	}
	createdAt := h.CatalogRepo.GetRecipeCreatedAt(ctx, recipeIDStr)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AiGenerationEvent{
		Source: "ai",
		Status: "complete",
		Recipe: toCompleteRecipe(detail, createdAt),
	})

	h.recordInteraction(ctx, userID, req.Ingredients, 1, &recipeIDStr, true)
}

// AiRecipeDetailHandler handles GET /api/ai/recipes/{id}.
func (h *AIChefHandler) AiRecipeDetailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/ai/recipes/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	detail, err := h.CatalogRepo.GetAiRecipeDetail(r.Context(), id)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	createdAt := h.CatalogRepo.GetRecipeCreatedAt(r.Context(), id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toCompleteRecipe(detail, createdAt))
}

// AiFeedbackHandler handles POST /api/ai/recipes/{id}/feedback.
func (h *AIChefHandler) AiFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/ai/recipes/")
	id = strings.TrimSuffix(id, "/feedback")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, "")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if body.Action != "save" && body.Action != "discard" {
		http.Error(w, "Action must be 'save' or 'discard'", http.StatusBadRequest)
		return
	}

	recipeID, _ := strconv.Atoi(id)
	userInt, _ := strconv.Atoi(userID)
	h.CatalogRepo.RecordFeedback(r.Context(), userInt, recipeID, body.Action)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// --- Helpers ---

func writeSSEEvent(w http.ResponseWriter, event string, data interface{}, flusher http.Flusher) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, string(jsonBytes))
	flusher.Flush()
}

func (h *AIChefHandler) sendError(w http.ResponseWriter, msg, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorEvent{Error: msg, Code: code})
}

func (h *AIChefHandler) sendErrorPlain(w http.ResponseWriter, msg, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorEvent{Error: msg, Code: code})
}

func (h *AIChefHandler) sendSSEError(w http.ResponseWriter, msg, code string, flusher http.Flusher) {
	writeSSEEvent(w, "error", ErrorEvent{Error: msg, Code: code}, flusher)
}

func (h *AIChefHandler) recordInteraction(ctx context.Context, userID string, ingredients []string, resultsCount int, selectedID *string, wasAI bool) {
	var selInt *int
	if selectedID != nil {
		id, _ := strconv.Atoi(*selectedID)
		selInt = &id
	}
	userInt, _ := strconv.Atoi(userID)
	h.CatalogRepo.RecordInteraction(ctx, userInt, ingredients, resultsCount, selInt, wasAI)
}

func toRecipeSummaries(rs []recipes.RecipeSummary) []RecipeSummary {
	out := make([]RecipeSummary, len(rs))
	for i, r := range rs {
		out[i] = RecipeSummary(r)
	}
	return out
}

func toCompleteRecipe(d *recipes.AiRecipeDetail, createdAt string) *AiCompleteRecipe {
	return &AiCompleteRecipe{
		ID:               d.ID,
		Name:             d.Name,
		Description:      d.Description,
		ImgURL:           d.ImgURL,
		Ingredients:      toIngredientOuts(d.Ingredients),
		Steps:            toStepOuts(d.Steps),
		Tags:             []string{},
		EstimatedMinutes: safeIntPtr(d.EstimatedMinutes),
		Difficulty:       safeStrPtr(d.Difficulty),
		Cuisine:          safeStrPtr(d.Cuisine),
		AiGenerated:      d.AiGenerated,
		CreatedAt:        createdAt,
	}
}

func toIngredientOuts(ings []recipes.AiIngredientOut) []AiIngredientOut {
	out := make([]AiIngredientOut, len(ings))
	for i, ing := range ings {
		out[i] = AiIngredientOut(ing)
	}
	return out
}

func toStepOuts(steps []recipes.AiStepOut) []AiStepOut {
	out := make([]AiStepOut, len(steps))
	for i, s := range steps {
		out[i] = AiStepOut(s)
	}
	return out
}

func toAiChefConstraints(c *AiConstraints) *ai.AiChefRequestConstraints {
	if c == nil {
		return nil
	}
	return &ai.AiChefRequestConstraints{
		MaxTime:    c.MaxTime,
		Difficulty: c.Difficulty,
		Cuisine:    c.Cuisine,
	}
}

func safeIntPtr(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func safeStrPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
