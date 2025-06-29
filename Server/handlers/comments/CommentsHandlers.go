package comments

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/utils"
)

func (h *CommentHandler) CommentsHandler(w http.ResponseWriter, r *http.Request)  {
	// Method check
	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query("SELECT c.id, u.name, c.recipe_id, c.comment, c.rating FROM comments c JOIN users u ON u.id = c.user_id WHERE recipe_id = $1", id) 	
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next(){
		var c Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.RecipeID, &c.Comment, &c.Rating); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) PostComment(w http.ResponseWriter, r *http.Request)  {
	// Method check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Recipe ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/comments/post/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	userID, err := utils.ParseToken(r, h.SecretKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Read JSON body
	var input struct {
		Comment string `json:"comment"`
		Rating float32 `json:"rating"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO comments (user_id, recipe_id, comment, rating) VALUES ($1, $2, $3, $4)",
		userID, id, input.Comment, input.Rating,
	)
	if err != nil {
		log.Println("DB error", err)
		http.Error(w, "Error creating a comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Comment created"))
}
