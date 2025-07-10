package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Zheng5005/BiteBox/db"
	"github.com/Zheng5005/BiteBox/handlers/auth"
	"github.com/Zheng5005/BiteBox/handlers/comments"
	"github.com/Zheng5005/BiteBox/handlers/meals"
	"github.com/Zheng5005/BiteBox/handlers/recipes"
	"github.com/Zheng5005/BiteBox/handlers/users"
	"github.com/Zheng5005/BiteBox/middlewares"
)

func main() {
	db.InitDB()
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "other_key"
	}

	commentHandler := comments.NewCommentHandler(db.DB, secret)
	recipesHandler := recipes.NewRecipesHandler(db.DB, secret)
	authHandler := auth.NewAuthHandler(db.DB, secret)
	userHandler := users.NewUserHandler(db.DB, secret)

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/signup", authHandler.SignUpHandler)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)

	// Users routes
	mux.HandleFunc("/api/users", middleware.JWTMiddleware(userHandler.GetRecipesAuth))
	mux.HandleFunc("PATCH /api/users/edit/", middleware.JWTMiddleware(userHandler.EditRecipeAuth))
	mux.HandleFunc("PATCH /api/users/deactivate/", middleware.JWTMiddleware(userHandler.DeActivateRecipeAuth))
	mux.HandleFunc("PATCH /api/users/activate/", middleware.JWTMiddleware(userHandler.ActivateRecipeAuth))

	mux.HandleFunc("GET /api/users/ByUser", userHandler.GetRecipesByUser)
	mux.HandleFunc("GET /api/users/ByGuest", userHandler.GetRecipesByGuestName)

	// Recipes routes
	mux.HandleFunc("/api/recipes", recipesHandler.RecipeHandler)
	mux.HandleFunc("/api/recipes/", recipesHandler.RecipeONEHandler)
	mux.HandleFunc("/api/recipes/guestPost", recipesHandler.PostRecipeGuest)
	mux.HandleFunc("/api/recipes/userPost", middleware.JWTMiddleware(recipesHandler.PostRecipeUser))

	// Comments routes
	mux.HandleFunc("/api/comments/", commentHandler.CommentsHandler)
	mux.HandleFunc("/api/comments/post/", middleware.JWTMiddleware(commentHandler.PostComment))

	// Meals routes
	mux.HandleFunc("/api/mealtypes", meals.MealsHandler)

	// CORS
	handlerWithCORS := middleware.CorsMiddleware(mux)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}

