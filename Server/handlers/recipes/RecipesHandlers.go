package recipes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (h *RecipesHandler) RecipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := h.DB.Query(`
		SELECT 
			r.id, 
			r.name_recipe, 
			r.description, 
			r.meal_type_id, 
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg 
		FROM recipes r 
		LEFT JOIN comments c ON r.id = c.recipe_id 
		GROUP BY r.id`)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []RecipesMainPage

	for rows.Next() {
		var r RecipesMainPage
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.MealTypeID, &r.ImgURL, &r.Rating); err != nil {
			log.Println(err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipesHandler) RecipeONEHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/recipes/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}
			
	query := `
			SELECT 
				r.id,
				r.name_recipe,
				r.description,
				r.meal_type_id,
				COALESCE(r.img_url, ''),
				COALESCE(u.name, r.guest_name) AS creator_name,
				COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg_rating,
				r.steps
			FROM recipes r
			LEFT JOIN users u ON u.id = r.user_id
			LEFT JOIN comments c ON c.recipe_id = r.id
			WHERE r.id = $1
			GROUP BY r.id, u.name, r.guest_name;
		`

		var recipe RecipeDetail
			
		err := h.DB.QueryRow(query, id).Scan(
			&recipe.ID,
			&recipe.Name,
			&recipe.Description,
			&recipe.MealTypeID,
			&recipe.ImgURL,
			&recipe.CreatorName,
			&recipe.Rating,
			&recipe.Steps,
		)

		if err == sql.ErrNoRows {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Println("Scan error:", err)
			http.Error(w, "Error retrieving recipe", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(recipe)
}

func (h *RecipesHandler) PostRecipeGuest(w http.ResponseWriter, r *http.Request)  {
	
}

func (h *RecipesHandler) PostRecipeUser(w http.ResponseWriter, r *http.Request)  {
	
}
