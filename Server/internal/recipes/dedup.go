package recipes

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/Zheng5005/BiteBox/internal/ai"
)

const dedupTitleWeight = 0.6
const dedupIngredientWeight = 0.4
const dedupThreshold = 0.8

// FindNearDuplicate checks if a near-duplicate of the candidate AI recipe already exists.
// Returns the existing recipe if combined similarity score > 0.8, otherwise nil.
func (r *CatalogRepo) FindNearDuplicate(ctx context.Context, candidate *ai.AiRecipe) (*RecipeSummary, error) {
	// Search among active recipes that share at least one tag
	query := `
		SELECT r.id::text, r.name_recipe, r.description,
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0)::text AS rating,
			COUNT(DISTINCT rl.id) AS likes,
			0 AS match_score
		FROM recipes r
		LEFT JOIN comments c ON r.id = c.recipe_id
		LEFT JOIN recipe_likes rl ON r.id = rl.recipe_id
		WHERE r.is_active = true
			AND r.ai_generated = true
			AND r.tags && $1::TEXT[]
		GROUP BY r.id
		ORDER BY r.id DESC
		LIMIT 50
	`

	candidateTags := pqArray(candidate.Tags)
	rows, err := r.DB.QueryContext(ctx, query, candidateTags)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("dedup query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bestMatch *RecipeSummary
	var bestScore float64

	for rows.Next() {
		var s RecipeSummary
		var likes sql.NullInt64
		var name, description string
		if err := rows.Scan(&s.ID, &name, &description, &s.ImgURL, &s.Rating, &likes, &s.MatchScore); err != nil {
			continue
		}
		if likes.Valid {
			s.Likes = int(likes.Int64)
		}
		s.Name = name
		s.Description = description

		score := r.combinedScore(candidate.Title, candidate.Ingredients, name, s.ID)
		if score > bestScore {
			bestScore = score
			bestMatch = &s
		}
	}

	if bestScore > dedupThreshold && bestMatch != nil {
		bestMatch.MatchScore = bestScore
		return bestMatch, nil
	}
	return nil, nil
}

// combinedScore computes 0.6 * title_similarity + 0.4 * ingredient_overlap.
func (r *CatalogRepo) combinedScore(candidateTitle string, candidateIngredients []ai.AiIngredient, existingTitle string, existingID string) float64 {
	titleSim := levenshteinRatio(candidateTitle, existingTitle)

	// Fetch existing ingredients for this recipe
	existingNames, err := r.getIngredientNames(existingID)
	if err != nil || len(existingNames) == 0 {
		return titleSim * dedupTitleWeight // fallback to title only
	}

	candidateNames := make([]string, len(candidateIngredients))
	for i, ing := range candidateIngredients {
		candidateNames[i] = strings.ToLower(ing.Name)
	}

	ingredientSim := jaccardSimilarity(candidateNames, existingNames)

	return dedupTitleWeight*titleSim + dedupIngredientWeight*ingredientSim
}

func (r *CatalogRepo) getIngredientNames(recipeID string) ([]string, error) {
	rows, err := r.DB.Query(
		"SELECT ingredient_name FROM recipe_ingredients WHERE recipe_id = $1", recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		names = append(names, strings.ToLower(name))
	}
	return names, nil
}

// levenshteinRatio returns 1.0 - (distance / max(len_a, len_b)).
func levenshteinRatio(a, b string) float64 {
	la, lb := len(a), len(b)
	if la == 0 && lb == 0 {
		return 1.0
	}
	if la == 0 || lb == 0 {
		return 0.0
	}

	// Simple Levenshtein distance
	maxLen := la
	if lb > maxLen {
		maxLen = lb
	}
	dist := levenshteinDistance(a, b)
	return 1.0 - float64(dist)/float64(maxLen)
}

func levenshteinDistance(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)

	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			dp[i][j] = min3(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+cost)
		}
	}
	return dp[la][lb]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// jaccardSimilarity computes |A ∩ B| / |A ∪ B| for two string sets.
func jaccardSimilarity(a, b []string) float64 {
	setA := make(map[string]bool)
	for _, s := range a {
		setA[strings.ToLower(strings.TrimSpace(s))] = true
	}
	setB := make(map[string]bool)
	for _, s := range b {
		setB[strings.ToLower(strings.TrimSpace(s))] = true
	}

	intersection := 0
	for k := range setA {
		if setB[k] {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0.0
	}
	return float64(intersection) / float64(union)
}
