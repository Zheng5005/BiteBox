package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zheng5005/BiteBox/db"
	"github.com/Zheng5005/BiteBox/handlers/aichef"
	"github.com/Zheng5005/BiteBox/handlers/auth"
	"github.com/Zheng5005/BiteBox/handlers/comments"
	"github.com/Zheng5005/BiteBox/handlers/cookbooks"
	"github.com/Zheng5005/BiteBox/handlers/meals"
	"github.com/Zheng5005/BiteBox/handlers/recipes"
	"github.com/Zheng5005/BiteBox/handlers/users"
	"github.com/Zheng5005/BiteBox/internal/ai"
	"github.com/Zheng5005/BiteBox/internal/ratelimit"
	intRecipes "github.com/Zheng5005/BiteBox/internal/recipes"
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
	cookbookHandler := cookbooks.NewCookbookHandler(db.DB, secret)

	// AI Chef setup
	aiClient := ai.NewAIClient(os.Getenv("GEMINI_API_KEY"), "gemini-2.0-flash")
	catalogRepo := intRecipes.NewCatalogRepo(db.DB)
	rateLimiter := ratelimit.NewRateLimiter(10, 24*time.Hour)
	rateLimiter.StartCleanup(1 * time.Hour)
	aiChefHandler := aichef.NewAIChefHandler(catalogRepo, aiClient, rateLimiter, 3)

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/signup", authHandler.SignUpHandler)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)

	// Users routes
	mux.HandleFunc("/api/users", middleware.JWTMiddleware(userHandler.GetRecipesAuth))
	mux.HandleFunc("PATCH /api/users/edit/", middleware.JWTMiddleware(userHandler.EditRecipeAuth))
	mux.HandleFunc("PATCH /api/users/deactivate/", middleware.JWTMiddleware(userHandler.DeActivateRecipeAuth))
	mux.HandleFunc("PATCH /api/users/activate/", middleware.JWTMiddleware(userHandler.ActivateRecipeAuth))
	mux.HandleFunc("POST /api/users/like/", middleware.JWTMiddleware(userHandler.LikeRecipe))

	mux.HandleFunc("GET /api/users/ByUser", userHandler.GetRecipesByUser)
	mux.HandleFunc("GET /api/users/ByGuest", userHandler.GetRecipesByGuestName)

	// Recipes routes
	mux.HandleFunc("/api/recipes", recipesHandler.RecipeHandler)
	mux.HandleFunc("/api/recipes/", recipesHandler.RecipeONEHandler)
	mux.HandleFunc("/api/recipes/post", recipesHandler.PostRecipe)

	// Comments routes
	mux.HandleFunc("/api/comments/", commentHandler.CommentsHandler)
	mux.HandleFunc("/api/comments/post/", middleware.JWTMiddleware(commentHandler.PostComment))

	// Cookbooks routes
	mux.HandleFunc("GET /api/cookbooks", middleware.JWTMiddleware(cookbookHandler.GetCookbooks))
	mux.HandleFunc("GET /api/cookbooks/", middleware.JWTMiddleware(cookbookHandler.GetCookbook))
	mux.HandleFunc("POST /api/cookbooks", middleware.JWTMiddleware(cookbookHandler.CreateCookbook))
	mux.HandleFunc("PATCH /api/cookbooks/edit/", middleware.JWTMiddleware(cookbookHandler.UpdateCookbook))
	mux.HandleFunc("DELETE /api/cookbooks/delete/", middleware.JWTMiddleware(cookbookHandler.DeleteCookbook))
	mux.HandleFunc("GET /api/cookbooks/recipes/", middleware.JWTMiddleware(cookbookHandler.GetCookbookRecipes))
	mux.HandleFunc("POST /api/cookbooks/recipes/add/", middleware.JWTMiddleware(cookbookHandler.AddRecipeToCookbook))
	mux.HandleFunc("DELETE /api/cookbooks/recipes/remove/", middleware.JWTMiddleware(cookbookHandler.RemoveRecipeFromCookbook))

	// Meals routes
	mux.HandleFunc("/api/mealtypes", meals.MealsHandler)

	// AI Chef routes
	mux.HandleFunc("POST /api/ai/recipes", middleware.JWTMiddleware(aiChefHandler.AiChefHandler))
	mux.HandleFunc("GET /api/ai/recipes/", aiChefHandler.AiRecipeDetailHandler)
	mux.HandleFunc("POST /api/ai/recipes/", middleware.JWTMiddleware(aiChefHandler.AiFeedbackHandler))

	// CORS
	handlerWithCORS := middleware.CorsMiddleware(mux)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}

