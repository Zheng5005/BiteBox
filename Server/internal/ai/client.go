package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models"

// AIClient handles communication with the Gemini API.
type AIClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewAIClient creates a new Gemini API client.
func NewAIClient(apiKey, model string) *AIClient {
	return &AIClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

// Generate sends a prompt to Gemini and returns a structured AiRecipe.
func (c *AIClient) Generate(ctx context.Context, prompt string) (*AiRecipe, error) {
	url := fmt.Sprintf("%s/%s:generateContent?key=%s", geminiAPIURL, c.model, c.apiKey)

	body := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiConfig{
			ResponseMimeType: "application/json",
			ResponseSchema:   recipeSchema(),
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error %d: %s", resp.StatusCode, string(respBody))
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from AI")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text
	return parseRecipeJSON(text)
}

// GenerateStream sends a prompt and returns a channel of streaming JSON tokens.
func (c *AIClient) GenerateStream(ctx context.Context, prompt string) (<-chan string, <-chan error) {
	tokenCh := make(chan string, 32)
	errCh := make(chan error, 1)

	url := fmt.Sprintf("%s/%s:streamGenerateContent?alt=sse&key=%s", geminiAPIURL, c.model, c.apiKey)

	body := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiConfig{
			ResponseMimeType: "application/json",
			ResponseSchema:   recipeSchema(),
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		errCh <- fmt.Errorf("marshal request: %w", err)
		close(tokenCh)
		close(errCh)
		return tokenCh, errCh
	}

	go func() {
		defer close(tokenCh)
		defer close(errCh)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			errCh <- fmt.Errorf("create request: %w", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("api call: %w", err)
			return
		}
		defer resp.Body.Close()

		buf := make([]byte, 4096)
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				tokenCh <- string(buf[:n])
			}
			if readErr != nil {
				if readErr == io.EOF {
					return
				}
				errCh <- fmt.Errorf("stream read: %w", readErr)
				return
			}
		}
	}()

	return tokenCh, errCh
}

// --- Gemini request/response types ---

type geminiRequest struct {
	Contents         []geminiContent `json:"contents"`
	GenerationConfig geminiConfig    `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiConfig struct {
	ResponseMimeType string           `json:"responseMimeType"`
	ResponseSchema   *recipeSchemaDef `json:"responseSchema,omitempty"`
}

type recipeSchemaDef struct {
	Type       string              `json:"type"`
	Properties map[string]*propDef `json:"properties"`
	Required   []string            `json:"required"`
}

type propDef struct {
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
	Items       *propDef            `json:"items,omitempty"`
	Properties  map[string]*propDef `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
	Format      string              `json:"format,omitempty"`
}

func recipeSchema() *recipeSchemaDef {
	return &recipeSchemaDef{
		Type: "object",
		Properties: map[string]*propDef{
			"title":             {Type: "string", Description: "Recipe title"},
			"description":       {Type: "string", Description: "1-2 sentence description"},
			"estimated_minutes": {Type: "integer", Format: "int32", Description: "Cooking time in minutes"},
			"difficulty":        {Type: "string", Enum: []string{"easy", "medium", "hard"}},
			"cuisine":           {Type: "string", Description: "Cuisine type"},
			"tags": {
				Type:  "array",
				Items: &propDef{Type: "string"},
			},
			"ingredients": {
				Type: "array",
				Items: &propDef{
					Type: "object",
					Properties: map[string]*propDef{
						"name":     {Type: "string"},
						"quantity": {Type: "number", Format: "double"},
						"unit":     {Type: "string"},
					},
					Required: []string{"name", "quantity", "unit"},
				},
			},
			"steps": {
				Type: "array",
				Items: &propDef{
					Type: "object",
					Properties: map[string]*propDef{
						"order":       {Type: "integer", Format: "int32"},
						"instruction": {Type: "string"},
					},
					Required: []string{"order", "instruction"},
				},
			},
		},
		Required: []string{"title", "description", "ingredients", "steps", "tags", "estimated_minutes", "difficulty", "cuisine"},
	}
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func parseRecipeJSON(raw string) (*AiRecipe, error) {
	var recipe AiRecipe
	if err := json.Unmarshal([]byte(raw), &recipe); err != nil {
		return nil, fmt.Errorf("parse ai recipe: %w", err)
	}
	if err := recipe.Validate(); err != nil {
		return nil, err
	}
	return &recipe, nil
}
