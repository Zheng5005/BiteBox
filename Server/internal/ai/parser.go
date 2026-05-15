package ai

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var jsonBlockRe = regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")

// ExtractJSON extracts JSON from raw AI output, handling markdown code fences.
func ExtractJSON(raw string) (string, error) {
	// Try markdown code block first
	matches := jsonBlockRe.FindStringSubmatch(raw)
	if len(matches) > 1 {
		return matches[1], nil
	}

	// Try to find JSON object boundaries
	start := -1
	for i := 0; i < len(raw); i++ {
		if raw[i] == '{' {
			start = i
			break
		}
	}
	if start == -1 {
		return "", fmt.Errorf("no JSON object found in response")
	}

	// Find matching closing brace
	depth := 0
	for i := start; i < len(raw); i++ {
		switch raw[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return raw[start : i+1], nil
			}
		}
	}

	return "", fmt.Errorf("unmatched JSON braces in response")
}

// ParseAndValidate extracts JSON from raw output and validates against AiRecipe schema.
func ParseAndValidate(raw string) (*AiRecipe, error) {
	jsonStr, err := ExtractJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("extract json: %w", err)
	}

	var recipe AiRecipe
	if err := json.Unmarshal([]byte(jsonStr), &recipe); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if err := recipe.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &recipe, nil
}

// AccumulateTokens reads from a token channel and returns the full accumulated response.
func AccumulateTokens(tokenCh <-chan string) string {
	var result string
	for token := range tokenCh {
		result += token
	}
	return result
}

// AccumulateAndParse reads from a token channel, accumulates, and parses the final result.
func AccumulateAndParse(tokenCh <-chan string) (*AiRecipe, error) {
	raw := AccumulateTokens(tokenCh)
	if raw == "" {
		return nil, fmt.Errorf("empty token stream")
	}
	return ParseAndValidate(raw)
}
