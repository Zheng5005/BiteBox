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
	ID          string `json:"id"`
	Name        string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID  string `json:"meal_type_id"`
	ImgURL      string `json:"img_url"`
	Rating      string `json:"rating"`
	Likes       int    `json:"likes"`
	Notes       string `json:"notes"`
}

type CookbookHandler struct {
	DB        db.DBExecutor
	SecretKey string
}

func NewCookbookHandler(db db.DBExecutor, secret string) *CookbookHandler {
	return &CookbookHandler{DB: db, SecretKey: secret}
}
