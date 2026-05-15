package recipes

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Zheng5005/BiteBox/db"
)

// SearchConstraints holds optional filters applied during catalog search.
type SearchConstraints struct {
	MaxTime    *int
	Difficulty *string
	Cuisine    *string
}

// SearchByIngredients finds recipes matching at least 2 of the provided ingredient tags.
func (r *CatalogRepo) SearchByIngredients(ctx context.Context, tags []string, constraints *SearchConstraints) ([]RecipeSummary, error) {
	whereExtras := ""
	args := []any{pqArray(tags)}
	argIdx := 2

	if constraints != nil {
		if constraints.MaxTime != nil && *constraints.MaxTime > 0 {
			whereExtras += fmt.Sprintf(" AND r.estimated_minutes <= $%d", argIdx)
			args = append(args, *constraints.MaxTime)
			argIdx++
		}
		if constraints.Difficulty != nil && *constraints.Difficulty != "" {
			whereExtras += fmt.Sprintf(" AND r.difficulty = $%d", argIdx)
			args = append(args, *constraints.Difficulty)
			argIdx++
		}
		if constraints.Cuisine != nil && *constraints.Cuisine != "" {
			whereExtras += fmt.Sprintf(" AND LOWER(r.cuisine) = LOWER($%d)", argIdx)
			args = append(args, *constraints.Cuisine)
			argIdx++
		}
	}

	query := fmt.Sprintf(`
		SELECT
			r.id::text,
			r.name_recipe,
			r.description,
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0)::text AS rating,
			COUNT(DISTINCT rl.id) AS likes,
			ROUND(
				CAST(array_length(ARRAY(
					SELECT unnest(r.tags)
					INTERSECT
					SELECT unnest($1::TEXT[])
				), 1) AS numeric)
				/ GREATEST(CAST(array_length($1::TEXT[], 1) AS numeric), 1),
				2
			) AS match_score
		FROM recipes r
		LEFT JOIN comments c ON r.id = c.recipe_id
		LEFT JOIN recipe_likes rl ON r.id = rl.recipe_id
		WHERE r.is_active = true
			AND r.tags && $1::TEXT[]
			AND array_length(ARRAY(
				SELECT unnest(r.tags)
				INTERSECT
				SELECT unnest($1::TEXT[])
			), 1) >= 2
			%s
		GROUP BY r.id
		ORDER BY match_score DESC, r.id DESC
		LIMIT 10
	`, whereExtras)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("catalog search error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []RecipeSummary
	for rows.Next() {
		var s RecipeSummary
		var likes sql.NullInt64
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.ImgURL, &s.Rating, &likes, &s.MatchScore); err != nil {
			log.Printf("catalog scan error: %v", err)
			return nil, err
		}
		if likes.Valid {
			s.Likes = int(likes.Int64)
		}
		results = append(results, s)
	}

	return results, nil
}

// pqArray converts a Go string slice to a PostgreSQL array literal.
func pqArray(arr []string) string {
	if len(arr) == 0 {
		return "{}"
	}
	result := "{"
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		s = strings.ReplaceAll(s, `\`, `\\`)
		s = strings.ReplaceAll(s, `"`, `\"`)
		result += `"` + s + `"`
	}
	result += "}"
	return result
}

// NewCatalogRepo creates a new CatalogRepo.
func NewCatalogRepo(db db.DBExecutor) *CatalogRepo {
	return &CatalogRepo{DB: db}
}
