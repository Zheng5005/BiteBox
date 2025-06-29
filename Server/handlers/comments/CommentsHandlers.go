package comments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid method")
		}
		return []byte(h.SecretKey), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Missing user ID in token", http.StatusUnauthorized)
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
