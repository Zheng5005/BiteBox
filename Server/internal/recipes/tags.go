package recipes

import (
	"strings"
)

// commonWords removed from ingredient names during normalization.
var commonWords = map[string]bool{
	"a": true, "an": true, "the": true, "some": true, "fresh": true,
	"of": true, "and": true, "or": true, "with": true, "to": true,
	"for": true, "in": true, "on": true, "by": true,
}

// NormalizeIngredients converts raw ingredient strings into canonical tags.
// Lowercases, trims, removes stop words, singularizes, and deduplicates.
func NormalizeIngredients(raw []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range raw {
		tag := normalizeTag(item)
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		result = append(result, tag)
	}

	return result
}

// normalizeTag processes a single ingredient string.
func normalizeTag(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))

	// Remove stop words
	words := strings.Fields(s)
	var filtered []string
	for _, w := range words {
		if !commonWords[w] {
			filtered = append(filtered, w)
		}
	}
	if len(filtered) == 0 {
		return ""
	}

	s = strings.Join(filtered, "_")

	// Basic singularization
	s = singularize(s)

	return s
}

// singularize applies simple English singularization rules.
func singularize(s string) string {
	if strings.HasSuffix(s, "ies") && len(s) > 4 {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "ves") && len(s) > 4 {
		return s[:len(s)-3] + "fe"
	}
	if strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "ches") || strings.HasSuffix(s, "shes") {
		return s[:len(s)-2]
	}
	if strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss") && len(s) > 2 {
		return s[:len(s)-1]
	}
	return s
}
