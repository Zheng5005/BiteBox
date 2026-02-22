package cookbooks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/utils"
)

func (h *CookbookHandler) GetCookbooks(w http.ResponseWriter, r *http.Request) {
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
		SELECT id, user_id, name, COALESCE(description, ''), is_public, created_at, updated_at
		FROM cookbooks
		WHERE user_id = $1`, userID)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cookbooks []Cookbook

	for rows.Next() {
		var c Cookbook
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Description, &c.IsPublic, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Println(err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		cookbooks = append(cookbooks, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cookbooks)
}

func (h *CookbookHandler) GetCookbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cookbooks/")
	if id == "" {
		http.Error(w, "Missing cookbook ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var c Cookbook
	err = h.DB.QueryRow(`
		SELECT id, user_id, name, COALESCE(description, ''), is_public, created_at, updated_at
		FROM cookbooks
		WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&c.ID, &c.UserID, &c.Name, &c.Description, &c.IsPublic, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Cookbook not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Scan error:", err)
		http.Error(w, "Error retrieving cookbook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (h *CookbookHandler) CreateCookbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		http.Error(w, "Missing required field: name", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO cookbooks (user_id, name, description, is_public) VALUES ($1, $2, $3, $4)",
		userID, body.Name, body.Description, body.IsPublic,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating cookbook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Cookbook Created"))
}

func (h *CookbookHandler) UpdateCookbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cookbooks/edit/")
	if id == "" {
		http.Error(w, "Missing cookbook ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		IsPublic    *bool   `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateFields := []string{}
	args := []interface{}{}
	i := 1

	if body.Name != nil {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", i))
		args = append(args, *body.Name)
		i++
	}
	if body.Description != nil {
		updateFields = append(updateFields, fmt.Sprintf("description = $%d", i))
		args = append(args, *body.Description)
		i++
	}
	if body.IsPublic != nil {
		updateFields = append(updateFields, fmt.Sprintf("is_public = $%d", i))
		args = append(args, *body.IsPublic)
		i++
	}

	if len(updateFields) == 0 {
		http.Error(w, "No valid fields to update", http.StatusBadRequest)
		return
	}

	updateFields = append(updateFields, "updated_at = CURRENT_TIMESTAMP")

	args = append(args, id, userID)
	query := fmt.Sprintf("UPDATE cookbooks SET %s WHERE id = $%d AND user_id = $%d",
		strings.Join(updateFields, ", "), i, i+1,
	)

	res, err := h.DB.Exec(query, args...)
	if err != nil {
		log.Println("DB update error:", err)
		http.Error(w, "Failed to update cookbook", http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Cookbook not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cookbook Updated"))
}

func (h *CookbookHandler) DeleteCookbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cookbooks/delete/")
	if id == "" {
		http.Error(w, "Missing cookbook ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	res, err := h.DB.Exec("DELETE FROM cookbooks WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting cookbook", http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Cookbook not found or not owned by user", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cookbook Deleted"))
}

func (h *CookbookHandler) GetCookbookRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cookbooks/recipes/")
	if id == "" {
		http.Error(w, "Missing cookbook ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Verify the cookbook belongs to the user
	var exists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM cookbooks WHERE id = $1 AND user_id = $2)", id, userID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Cookbook not found", http.StatusNotFound)
		return
	}

	rows, err := h.DB.Query(`
		SELECT 
			r.id,
			r.name_recipe,
			r.description,
			r.meal_type_id,
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS rating,
			COUNT(DISTINCT rl.id) AS likes
		FROM cookbook_recipes cr
		JOIN recipes r ON r.id = cr.recipe_id
		LEFT JOIN comments c ON c.recipe_id = r.id
		LEFT JOIN recipe_likes rl ON rl.recipe_id = r.id
		WHERE cr.cookbook_id = $1 AND r.is_active = true
		GROUP BY r.id`, id)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []CookbookRecipe

	for rows.Next() {
		var cr CookbookRecipe
		if err := rows.Scan(&cr.ID, &cr.Name, &cr.Description, &cr.MealTypeID, &cr.ImgURL, &cr.Rating, &cr.Likes); err != nil {
			log.Println(err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, cr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func (h *CookbookHandler) AddRecipeToCookbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cookbooks/recipes/add/")
	if id == "" {
		http.Error(w, "Missing cookbook ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Verify the cookbook belongs to the user
	var exists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM cookbooks WHERE id = $1 AND user_id = $2)", id, userID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Cookbook not found", http.StatusNotFound)
		return
	}

	var body struct {
		RecipeID int    `json:"recipe_id"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.RecipeID == 0 {
		http.Error(w, "Missing required field: recipe_id", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO cookbook_recipes (cookbook_id, recipe_id, notes) VALUES ($1, $2, $3)",
		id, body.RecipeID, body.Notes,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding recipe to cookbook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Recipe added to cookbook"))
}

func (h *CookbookHandler) RemoveRecipeFromCookbook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/cookbooks/recipes/remove/")
	if id == "" {
		http.Error(w, "Missing cookbook ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Verify the cookbook belongs to the user
	var owns bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM cookbooks WHERE id = $1 AND user_id = $2)", id, userID).Scan(&owns)
	if err != nil || !owns {
		http.Error(w, "Cookbook not found", http.StatusNotFound)
		return
	}

	var body struct {
		RecipeID int `json:"recipe_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.RecipeID == 0 {
		http.Error(w, "Missing required field: recipe_id", http.StatusBadRequest)
		return
	}

	res, err := h.DB.Exec("DELETE FROM cookbook_recipes WHERE cookbook_id = $1 AND recipe_id = $2", id, body.RecipeID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error removing recipe from cookbook", http.StatusInternalServerError)
		return
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		http.Error(w, "Recipe not found in cookbook", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recipe removed from cookbook"))
}
