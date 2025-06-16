package recipes

import (
	"encoding/json"
	"net/http"

	"github.com/Zheng5005/BiteBox/db"
)

type Recipe struct {
	ID   string `json:"id"`
	UserID string `json:"user_id"`
	Name string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID string `json:"meal_type_id"`
	ImgURL string `json:"img_url"`
	GuestName string `json:"guest_name"`

	Rating string `json:"rating"`
}

type RecipesMainPage struct {
	ID   string `json:"id"`
	Name string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID string `json:"meal_type_id"`
	ImgURL string `json:"img_url"`

	Rating string `json:"rating"`
}

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			rows, err := db.DB.Query("SELECT r.id, r.name_recipe, r.description, r.meal_type_id, AVG(c.rating) FROM recipes r JOIN comments c ON r.id = c.recipe_id GROUP BY r.id")
			if err != nil {
				http.Error(w, "Query error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var recipes []RecipesMainPage

			for rows.Next() {
				var r RecipesMainPage
				if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.MealTypeID, &r.Rating); err != nil {
					http.Error(w, "Scan error", http.StatusInternalServerError)
					return
				}
				recipes = append(recipes, r)
			}

			json.NewEncoder(w).Encode(recipes)
		case http.MethodPost:
			
	}
}

