package ai

import (
	"fmt"
	"strings"
)

const systemPrompt = `You are a professional chef assistant. Generate a recipe based on the provided ingredients and constraints.

CRITICAL RULES:
1. Output ONLY valid JSON — no markdown code fences, no prose, no explanations.
2. Use only the ingredients provided by the user (do not invent extras).
3. Quantities must be numbers (e.g. 1.5, not "one and a half").
4. Units must be standard measurements (cups, tbsp, tsp, grams, ml, oz, pieces, etc.).
5. Steps must be numbered sequentially starting from 1.
6. Tags should include the main ingredients and cuisine type.
7. Difficulty must be one of: "easy", "medium", "hard".
8. Estimated minutes must be a realistic cooking time.`

// BuildPrompt creates the full prompt for AI recipe generation.
func BuildPrompt(ingredients []string, constraints *AiChefRequestConstraints, catalogHints []string) string {
	var sb strings.Builder

	sb.WriteString(systemPrompt)
	sb.WriteString("\n\n")

	sb.WriteString("User's available ingredients:\n")
	for _, ing := range ingredients {
		sb.WriteString(fmt.Sprintf("- %s\n", ing))
	}
	sb.WriteString("\n")

	if constraints != nil {
		sb.WriteString("Constraints:\n")
		if constraints.MaxTime != nil && *constraints.MaxTime > 0 {
			sb.WriteString(fmt.Sprintf("- Maximum cooking time: %d minutes\n", *constraints.MaxTime))
		}
		if constraints.Difficulty != nil && *constraints.Difficulty != "" {
			sb.WriteString(fmt.Sprintf("- Difficulty level: %s\n", *constraints.Difficulty))
		}
		if constraints.Cuisine != nil && *constraints.Cuisine != "" {
			sb.WriteString(fmt.Sprintf("- Cuisine style: %s\n", *constraints.Cuisine))
		}
		sb.WriteString("\n")
	}

	if len(catalogHints) > 0 {
		sb.WriteString("Similar recipes already exist (create something different):\n")
		for _, hint := range catalogHints {
			sb.WriteString(fmt.Sprintf("- %s\n", hint))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Return a JSON object with this exact structure:\n")
	sb.WriteString(`{"title":"...","description":"...","ingredients":[{"name":"...","quantity":0.0,"unit":"..."}],"steps":[{"order":1,"instruction":"..."}],"tags":["..."],"estimated_minutes":0,"difficulty":"easy|medium|hard","cuisine":"..."}`)

	return sb.String()
}

// AiChefRequestConstraints mirrors the handler constraints type.
// Defined here to avoid circular imports.
type AiChefRequestConstraints struct {
	MaxTime    *int
	Difficulty *string
	Cuisine    *string
}
