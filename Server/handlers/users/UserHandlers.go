package users

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Zheng5005/BiteBox/db"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	URLPhoto string `json:"url_photo"`
	GoogleID string `json:"google_id"`
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query("SELECT name, email FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Name, &u.Email); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
