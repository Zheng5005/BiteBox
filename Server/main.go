package main

import (
	"log"
	"net/http"
	"github.com/Zheng5005/BiteBox/db"
	"github.com/Zheng5005/BiteBox/handlers/users"
	"github.com/Zheng5005/BiteBox/handlers/recipes"
	"github.com/Zheng5005/BiteBox/middlewares"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()

	// Users routes
	mux.HandleFunc("/api/users", users.GetUsers)

	// Recipes routes
	mux.HandleFunc("/api/recipes", recipes.GetRecipes)

	// CORS
	handlerWithCORS := middleware.CorsMiddleware(mux)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}

