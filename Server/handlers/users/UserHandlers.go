package users

import (
	"encoding/json"
	"net/http"

	"github.com/Zheng5005/BiteBox/db"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

