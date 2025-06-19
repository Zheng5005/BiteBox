package meals

import (
	"encoding/json"
	"net/http"

	"github.com/Zheng5005/BiteBox/db"
)

type MealsType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func MealsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name FROM meal_type")
	if err != nil {
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var meals []MealsType

	for rows.Next() {
		var m MealsType
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		meals = append(meals, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}
