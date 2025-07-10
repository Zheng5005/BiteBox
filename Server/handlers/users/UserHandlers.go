package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/lib"
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
	
	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name_recipe")
	description := r.FormValue("description")
	mealType := r.FormValue("meal_type_id")
	steps := r.FormValue("steps")

	//Upload image to Cloudinary
	file, fileHeader, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		// Assume you have a cloudinary uploader function
		imageURL, err = lib.UploadToCloudinary(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Error uploading image", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	// Dynamic SQL
	updateFields := []string{}
	args := []interface{}{}
	i := 1

	if name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name_recipe = $%d", i))
    args = append(args, name)
    i++
  }
  if description != "" {
    updateFields = append(updateFields, fmt.Sprintf("description = $%d", i))
    args = append(args, description)
    i++
  }
  if mealType != "" {
    updateFields = append(updateFields, fmt.Sprintf("meal_type_id = $%d", i))
    args = append(args, mealType)
    i++
	}
	if steps != "" {
    updateFields = append(updateFields, fmt.Sprintf("steps = $%d", i))
    args = append(args, steps)
    i++
  }
  if imageURL != "" {
    updateFields = append(updateFields, fmt.Sprintf("img_url = $%d", i))
    args = append(args, imageURL)
    i++
  }

	if len(updateFields) == 0 {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
    return
	}

	args = append(args, id, userID)
	query := fmt.Sprintf("UPDATE recipes SET %s WHERE id = $%d AND user_id = $%d",
		strings.Join(updateFields, ", "),
    i, i+1,
  )

	res, err := h.DB.Exec(query,args...)
	if err != nil {
		log.Println("DB update error:", err)
    http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
    return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Recipe not found or not owned by user", http.StatusNotFound)
    return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recipe Updated"))
}

func (h *UserHandler) GetRecipesByUser(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user_name := r.URL.Query().Get("userName")

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
		LEFT JOIN users u ON r.user_id = u.id
		WHERE u.name = $1
		GROUP BY r.id`, user_name)
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

func (h *UserHandler) GetRecipesByGuestName(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guest_name := r.URL.Query().Get("guestName")

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
		WHERE r.guest_name = $1
		GROUP BY r.id`, guest_name)
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
