package recipes

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Zheng5005/BiteBox/internal/ai"
)

// SaveAiRecipe persists an AI-generated recipe into the database.
// Returns the new recipe ID.
func (r *CatalogRepo) SaveAiRecipe(ctx context.Context, recipe *ai.AiRecipe, sourcePrompt string, tags []string) (int, error) {
	stepsJSON, err := json.Marshal(recipe.Steps)
	if err != nil {
		return 0, err
	}

	var recipeID int
	err = r.DB.QueryRowContext(ctx, `
		INSERT INTO recipes (
			name_recipe, description, ai_generated, source_prompt,
			generated_at, estimated_minutes, difficulty, cuisine, tags, steps_json
		) VALUES ($1, $2, true, $3, NOW(), $4, $5, $6, $7, $8)
		RETURNING id
	`, recipe.Title, recipe.Description, sourcePrompt,
		recipe.EstimatedMinutes, recipe.Difficulty, recipe.Cuisine,
		pqArray(tags), stepsJSON).Scan(&recipeID)

	if err != nil {
		log.Printf("save ai recipe error: %v", err)
		return 0, err
	}

	// Insert structured ingredients
	for _, ing := range recipe.Ingredients {
		_, err := r.DB.ExecContext(ctx, `
			INSERT INTO recipe_ingredients (recipe_id, ingredient_name, quantity, unit)
			VALUES ($1, $2, $3, $4)
		`, recipeID, ing.Name, ing.Quantity, ing.Unit)
		if err != nil {
			log.Printf("save ingredient error: %v", err)
			// Non-fatal: recipe was saved, log and continue
		}
	}

	return recipeID, nil
}

// GetAiRecipeDetail retrieves a recipe with its structured ingredients.
func (r *CatalogRepo) GetAiRecipeDetail(ctx context.Context, recipeID string) (*AiRecipeDetail, error) {
	query := `
		SELECT
			r.id::text, r.name_recipe, r.description, r.meal_type_id,
			COALESCE(r.img_url, ''),
			COALESCE(u.name, r.guest_name) AS creator_name,
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0)::text AS avg_rating,
			COUNT(DISTINCT rl.id) AS likes,
			r.ai_generated, r.estimated_minutes, r.difficulty, r.cuisine,
			r.steps_json
		FROM recipes r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN comments c ON c.recipe_id = r.id
		LEFT JOIN recipe_likes rl ON rl.recipe_id = r.id
		WHERE r.id = $1 AND r.is_active = true
		GROUP BY r.id, u.name, r.guest_name
	`

	var detail AiRecipeDetail
	var mealTypeID sqlNullInt
	var stepsJSON []byte
	var aiGenerated bool

	err := r.DB.QueryRowContext(ctx, query, recipeID).Scan(
		&detail.ID, &detail.Name, &detail.Description, &mealTypeID,
		&detail.ImgURL, &detail.CreatorName, &detail.Rating, &detail.Likes,
		&aiGenerated, &detail.EstimatedMinutes, &detail.Difficulty, &detail.Cuisine,
		&stepsJSON,
	)
	if err != nil {
		return nil, err
	}

	if mealTypeID.Valid {
		mtid := int(mealTypeID.Int64)
		detail.MealTypeID = &mtid
	}

	// Parse structured steps from JSONB
	if len(stepsJSON) > 0 {
		var steps []AiStepOut
		if err := json.Unmarshal(stepsJSON, &steps); err == nil {
			detail.Steps = steps
		}
	}

	// Fetch structured ingredients
	ingredients, err := r.getIngredientsByRecipeID(recipeID)
	if err == nil {
		detail.Ingredients = ingredients
	}

	detail.AiGenerated = aiGenerated
	return &detail, nil
}

// AiStepOut represents a structured cooking step.
type AiStepOut struct {
	Order       int    `json:"order"`
	Instruction string `json:"instruction"`
}

// AiIngredientOut represents a structured ingredient.
type AiIngredientOut struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// AiRecipeDetail is the full recipe detail response for AI-generated recipes.
type AiRecipeDetail struct {
	ID               string            `json:"id"`
	Name             string            `json:"name_recipe"`
	Description      string            `json:"description"`
	MealTypeID       *int              `json:"meal_type_id,omitempty"`
	ImgURL           string            `json:"img_url"`
	CreatorName      string            `json:"creator_name"`
	Rating           string            `json:"rating"`
	Likes            int               `json:"likes"`
	AiGenerated      bool              `json:"ai_generated"`
	EstimatedMinutes *int              `json:"estimated_minutes,omitempty"`
	Difficulty       *string           `json:"difficulty,omitempty"`
	Cuisine          *string           `json:"cuisine,omitempty"`
	Steps            []AiStepOut       `json:"steps,omitempty"`
	Ingredients      []AiIngredientOut `json:"ingredients,omitempty"`
}

type sqlNullInt struct {
	Int64 int64
	Valid bool
}

func (s *sqlNullInt) Scan(value interface{}) error {
	if value == nil {
		s.Valid = false
		return nil
	}
	switch v := value.(type) {
	case int64:
		s.Int64 = v
		s.Valid = true
	case float64:
		s.Int64 = int64(v)
		s.Valid = true
	case nil:
		s.Valid = false
	default:
		s.Valid = false
	}
	return nil
}

func (r *CatalogRepo) getIngredientsByRecipeID(recipeID string) ([]AiIngredientOut, error) {
	rows, err := r.DB.QueryContext(context.Background(), `
		SELECT ingredient_name, COALESCE(quantity, 0), COALESCE(unit, '')
		FROM recipe_ingredients WHERE recipe_id = $1 ORDER BY id
	`, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []AiIngredientOut
	for rows.Next() {
		var ing AiIngredientOut
		if err := rows.Scan(&ing.Name, &ing.Quantity, &ing.Unit); err != nil {
			continue
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients, nil
}

// RecordInteraction logs an AI Chef interaction event.
func (r *CatalogRepo) RecordInteraction(ctx context.Context, userID int, ingredients []string, resultsCount int, selectedRecipeID *int, wasAIGenerated bool) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO ai_chef_interactions (user_id, ingredients_submitted, results_count, recipe_selected_id, was_ai_generated)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, pqArray(ingredients), resultsCount, selectedRecipeID, wasAIGenerated)
	return err
}

// RecordFeedback logs a save/discard action on an AI recipe.
func (r *CatalogRepo) RecordFeedback(ctx context.Context, userID int, recipeID int, action string) error {
	var selectedID *int
	if action == "save" {
		selectedID = &recipeID
	}
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO ai_chef_interactions (user_id, recipe_selected_id, was_ai_generated, ingredients_submitted)
		VALUES ($1, $2, true, $3)
	`, userID, selectedID, pqArray([]string{"feedback:" + action}))
	return err
}

// FormatTime formats the recipe's creation time for API responses.
func (r *CatalogRepo) GetRecipeCreatedAt(ctx context.Context, recipeID string) string {
	var createdAt time.Time
	err := r.DB.QueryRowContext(ctx, "SELECT created_at FROM recipes WHERE id = $1", recipeID).Scan(&createdAt)
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	return createdAt.Format(time.RFC3339)
}
