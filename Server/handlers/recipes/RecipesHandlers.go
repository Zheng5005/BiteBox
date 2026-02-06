package recipes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/lib"
	"github.com/Zheng5005/BiteBox/utils"
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
		WHERE r.is_active = true
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
			WHERE r.id = $1 AND r.is_active = true
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name_recipe := r.FormValue("name")
	description := r.FormValue("description")
	steps := r.FormValue("steps")
	meal_type_id := r.FormValue("meal_type_id")
	guest_name := r.FormValue("guest_name")

	if name_recipe == "" || description == "" || steps == "" || meal_type_id == "" || guest_name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		imageURL, err = lib.UploadToCloudinary(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Error uploading image", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO recipes (guest_name, name_recipe, description, meal_type_id, img_url, steps) VALUES ($1, $2, $3, $4, $5, $6)",
		guest_name, name_recipe, description, meal_type_id, imageURL, steps,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Recipe Created"))
}

func (h *RecipesHandler) PostRecipeUser(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name_recipe := r.FormValue("name")
	description := r.FormValue("description")
	steps := r.FormValue("steps")
	meal_type_id := r.FormValue("meal_type_id")

	if name_recipe == "" || description == "" || steps == "" || meal_type_id == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		imageURL, err = lib.UploadToCloudinary(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Error uploading image", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO recipes (user_id, name_recipe, description, meal_type_id, img_url, steps) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, name_recipe, description, meal_type_id, imageURL, steps,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Recipe Created"))
}
