package recipes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"strconv"

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
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg,
			COUNT(DISTINCT rl.id) AS likes
		FROM recipes r 
		LEFT JOIN comments c ON r.id = c.recipe_id 
		LEFT JOIN recipe_likes rl ON r.id = rl.recipe_id
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
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.MealTypeID, &r.ImgURL, &r.Rating, &r.Likes); err != nil {
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

	// Increment views in background
	userID, errToken := utils.ParseToken(r, h.SecretKey)
	var userIDPtr *int
	if errToken == nil {
		userIDtoInt, err := strconv.Atoi(userID)
		if err != nil {
			log.Println("Parse error", err)
		}
		userIDPtr = &userIDtoInt
	}
	go h.incrementViews(id, userIDPtr)
			
	query := `
			SELECT 
				r.id,
				r.name_recipe,
				r.description,
				r.meal_type_id,
				COALESCE(r.img_url, ''),
				COALESCE(u.name, r.guest_name) AS creator_name,
				COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg_rating,
				r.steps,
				COUNT(DISTINCT rl.id) AS likes
			FROM recipes r
			LEFT JOIN users u ON u.id = r.user_id
			LEFT JOIN comments c ON c.recipe_id = r.id
			LEFT JOIN recipe_likes rl ON rl.recipe_id = r.id
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
			&recipe.Likes,
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

func (h *RecipesHandler) incrementViews(recipeID string, userID *int) {
	var err error
	if userID != nil {
		_, err = h.DB.Exec("INSERT INTO recipe_views (recipe_id, user_id) VALUES ($1, $2)", recipeID, *userID)
	} else {
		_, err = h.DB.Exec("INSERT INTO recipe_views (recipe_id) VALUES ($1)", recipeID)
	}
	if err != nil {
		log.Printf("Error incrementing views for recipe %s: %v", recipeID, err)
	}
}

func (h *RecipesHandler) GetPopularRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Score = (Views * 1) + (Likes * 5) + (AvgRating * 10)
	query := `
		SELECT 
			r.id, 
			r.name_recipe, 
			r.description, 
			r.meal_type_id, 
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg,
			COUNT(DISTINCT rl.id) AS likes,
			(
				COUNT(DISTINCT rv.id) * 1 + 
				COUNT(DISTINCT rl.id) * 5 + 
				COALESCE(AVG(c.rating), 0) * 10
			) AS score
		FROM recipes r 
		LEFT JOIN comments c ON r.id = c.recipe_id 
		LEFT JOIN recipe_likes rl ON r.id = rl.recipe_id
		LEFT JOIN recipe_views rv ON r.id = rv.recipe_id
		WHERE r.is_active = true
		GROUP BY r.id
		ORDER BY score DESC
		LIMIT 10`

	rows, err := h.DB.Query(query)
	if err != nil {
		log.Println("Popular recipes query error:", err)
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []RecipesMainPage
	for rows.Next() {
		var r RecipesMainPage
		var score float64
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.MealTypeID, &r.ImgURL, &r.Rating, &r.Likes, &score); err != nil {
			log.Println("Scan error:", err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipesHandler) GetRecommendedRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, errToken := utils.ParseToken(r, h.SecretKey)
	if errToken != nil {
		// Fallback for non-authenticated: just popular
		h.GetPopularRecipes(w, r)
		return
	}

	// 1. Get user's preferred meal types
	prefQuery := `
		SELECT meal_type_id FROM (
			SELECT r.meal_type_id FROM recipes r JOIN recipe_views rv ON r.id = rv.recipe_id WHERE rv.user_id = $1
			UNION ALL
			SELECT r.meal_type_id FROM recipes r JOIN recipe_likes rl ON r.id = rl.recipe_id WHERE rl.user_id = $1
		) AS interactions
		GROUP BY meal_type_id
		ORDER BY COUNT(*) DESC
		LIMIT 3`

	rowsPref, err := h.DB.Query(prefQuery, userID)
	if err != nil {
		log.Println("Preference query error:", err)
		h.GetPopularRecipes(w, r)
		return
	}
	defer rowsPref.Close()

	var mealTypeIDs []int
	for rowsPref.Next() {
		var mtID int
		if err := rowsPref.Scan(&mtID); err == nil {
			mealTypeIDs = append(mealTypeIDs, mtID)
		}
	}

	if len(mealTypeIDs) == 0 {
		h.GetPopularRecipes(w, r)
		return
	}

	// Build placeholders for the IN clause (e.g., $1, $2, $3)
	placeholders := make([]string, len(mealTypeIDs))
	args := make([]interface{}, len(mealTypeIDs))
	for i, id := range mealTypeIDs {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		args[i] = id
	}

	query := `
		SELECT 
			r.id, 
			r.name_recipe, 
			r.description, 
			r.meal_type_id, 
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg,
			COUNT(DISTINCT rl.id) AS likes,
			(
				COUNT(DISTINCT rv.id) * 1 + 
				COUNT(DISTINCT rl.id) * 5 + 
				COALESCE(AVG(c.rating), 0) * 10
			) AS score
		FROM recipes r 
		LEFT JOIN comments c ON r.id = c.recipe_id 
		LEFT JOIN recipe_likes rl ON r.id = rl.recipe_id
		LEFT JOIN recipe_views rv ON r.id = rv.recipe_id
		WHERE r.is_active = true AND r.meal_type_id IN (` + strings.Join(placeholders, ",") + `)
		GROUP BY r.id
		ORDER BY score DESC
		LIMIT 10`

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		log.Println("Recommended recipes query error:", err)
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []RecipesMainPage
	for rows.Next() {
		var r RecipesMainPage
		var score float64
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.MealTypeID, &r.ImgURL, &r.Rating, &r.Likes, &score); err != nil {
			log.Println("Scan error:", err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *RecipesHandler) PostRecipe(w http.ResponseWriter, r *http.Request) {
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

	userID, tokenErr := utils.ParseToken(r, h.SecretKey)

	if tokenErr == nil {
		_, err = h.DB.Exec(
			"INSERT INTO recipes (user_id, name_recipe, description, meal_type_id, img_url, steps) VALUES ($1, $2, $3, $4, $5, $6)",
			userID, name_recipe, description, meal_type_id, imageURL, steps,
		)
	} else {
		guest_name := r.FormValue("guest_name")
		if guest_name == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		_, err = h.DB.Exec(
			"INSERT INTO recipes (guest_name, name_recipe, description, meal_type_id, img_url, steps) VALUES ($1, $2, $3, $4, $5, $6)",
			guest_name, name_recipe, description, meal_type_id, imageURL, steps,
		)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating recipe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Recipe Created"))
}
