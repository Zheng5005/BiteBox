package main

import (
	"log"
	"net/http"

	"github.com/Zheng5005/BiteBox/db"
	"github.com/Zheng5005/BiteBox/handlers/auth"
	"github.com/Zheng5005/BiteBox/handlers/meals"
	"github.com/Zheng5005/BiteBox/handlers/recipes"
	"github.com/Zheng5005/BiteBox/handlers/comments"
	"github.com/Zheng5005/BiteBox/handlers/users"
	"github.com/Zheng5005/BiteBox/middlewares"
)

func main() {
	db.InitDB()

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/signup", auth.SignUpHandler)
	mux.HandleFunc("/api/auth/login", auth.LoginHandler)

	// Users routes
	mux.HandleFunc("/api/users/", middleware.JWTMiddleware(users.GetUserInfo))

	// Recipes routes
	mux.HandleFunc("/api/recipes", recipes.RecipeHandler)
	mux.HandleFunc("/api/recipes/", recipes.RecipeONEHandler)

	// Comments routes
	mux.HandleFunc("/api/comments/", comments.CommentsHandler)
	mux.HandleFunc("/api/comments/post/", middleware.JWTMiddleware(comments.PostComment))

	// Meals routes
	mux.HandleFunc("/api/mealtypes", meals.MealsHandler)

	// CORS
	handlerWithCORS := middleware.CorsMiddleware(mux)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}

