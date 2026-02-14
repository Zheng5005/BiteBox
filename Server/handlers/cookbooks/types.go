package cookbooks

import "github.com/Zheng5005/BiteBox/db"

type Cookbook struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CookbookRecipe struct {
	ID       int    `json:"id"`
	RecipeID int    `json:"recipe_id"`
	AddedAt  string `json:"added_at"`
	Notes    string `json:"notes"`
}

type CookbookHandler struct {
	DB        db.DBExecutor
	SecretKey string
}

func NewCookbookHandler(db db.DBExecutor, secret string) *CookbookHandler {
	return &CookbookHandler{DB: db, SecretKey: secret}
}
