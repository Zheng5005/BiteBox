package aichef

// AiChefRequest is the incoming request body for AI Chef.
type AiChefRequest struct {
	Ingredients []string       `json:"ingredients"`
	Constraints *AiConstraints `json:"constraints,omitempty"`
}

// AiConstraints holds optional cooking constraints.
type AiConstraints struct {
	MaxTime    *int    `json:"max_time,omitempty"`
	Difficulty *string `json:"difficulty,omitempty"`
	Cuisine    *string `json:"cuisine,omitempty"`
}

// CatalogResultEvent is the SSE event for catalog search results.
type CatalogResultEvent struct {
	Source     string          `json:"source"`
	Recipes    []RecipeSummary `json:"recipes"`
	MatchCount int             `json:"match_count"`
}

// RecipeSummary mirrors internal/recipes.RecipeSummary for SSE serialization.
type RecipeSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name_recipe"`
	Description string  `json:"description"`
	ImgURL      string  `json:"img_url"`
	Rating      string  `json:"rating"`
	Likes       int     `json:"likes"`
	MatchScore  float64 `json:"match_score"`
}

// AiGenerationEvent is the SSE event for AI generation (partial or complete).
type AiGenerationEvent struct {
	Source string            `json:"source"`
	Status string            `json:"status"` // "starting" | "partial" | "complete"
	Token  string            `json:"token,omitempty"`
	Recipe *AiCompleteRecipe `json:"recipe,omitempty"`
}

// AiCompleteRecipe is the full recipe returned on AI generation completion.
type AiCompleteRecipe struct {
	ID               string            `json:"id"`
	Name             string            `json:"name_recipe"`
	Description      string            `json:"description"`
	ImgURL           string            `json:"img_url"`
	Ingredients      []AiIngredientOut `json:"ingredients"`
	Steps            []AiStepOut       `json:"steps"`
	Tags             []string          `json:"tags"`
	EstimatedMinutes int               `json:"estimated_minutes"`
	Difficulty       string            `json:"difficulty"`
	Cuisine          string            `json:"cuisine"`
	AiGenerated      bool              `json:"ai_generated"`
	CreatedAt        string            `json:"created_at"`
}

// AiIngredientOut is the ingredient format for API responses.
type AiIngredientOut struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// AiStepOut is the step format for API responses.
type AiStepOut struct {
	Order       int    `json:"order"`
	Instruction string `json:"instruction"`
}

// ErrorEvent is the SSE event for errors.
type ErrorEvent struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
