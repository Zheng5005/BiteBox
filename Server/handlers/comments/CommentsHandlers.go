package comments

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/db"
)

type Comment struct {
	ID   string `json:"id"`
	UserID string `json:"user_name"`
	RecipeID string `json:"recipe_id"`
	Comment string `json:"comment"`
	Rating string `json:"rating"`
}

func CommentsHandler(w http.ResponseWriter, r *http.Request)  {
	id := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	if id == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := db.DB.Query("SELECT c.id, u.name, c.recipe_id, c.comment, c.rating FROM comments c JOIN users u ON u.id = c.user_id WHERE recipe_id = $1", id) 	
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
	case http.MethodPost:
		//logic for POST (auth)
	}
}
