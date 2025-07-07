package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/utils"
)

func (h *UserHandler) GetRecipesAuth (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
		WHERE r.user_id = $1
		GROUP BY r.id`, userID)
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

func (h *UserHandler) DeActivateRecipeAuth(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/users/deactivate/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}
	
	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	query := `
			UPDATE recipes
			SET is_active = false
			WHERE id = $1 AND user_id = $2
		`

	_, err = h.DB.Exec(query, id, userID)
	if err != nil {
		log.Panicln(err)
		http.Error(w, "Error deactivating recipe", http.StatusNotModified)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recipe deactivated"))
}

func (h *UserHandler) ActivateRecipeAuth(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/users/activate/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}
	
	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	query := `
			UPDATE recipes
			SET is_active = true
			WHERE id = $1 AND user_id = $2
		`

	_, err = h.DB.Exec(query, id, userID)
	if err != nil {
		http.Error(w, "Error acativating recipe", http.StatusNotModified)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recipe activated"))
}

func (h *UserHandler) EditRecipeAuth(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/users/edit/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}
	
	_, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
}

func (h *UserHandler) GetRecipesByUser(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
}

func (h *UserHandler) GetRecipesByGuestName(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
}
