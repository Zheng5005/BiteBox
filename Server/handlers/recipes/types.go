package recipes

import "github.com/Zheng5005/BiteBox/db"

// This type is for POST and PATCH?
type Recipe struct {
	ID   string `json:"id"`
	UserID string `json:"user_id"`
	Name string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID string `json:"meal_type_id"`
	ImgURL string `json:"img_url"`
	GuestName string `json:"guest_name"`

	Rating string `json:"rating"`
}

//Type crafted with the main page in mind
type RecipesMainPage struct {
	ID   string `json:"id"`
	Name string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID string `json:"meal_type_id"`
	ImgURL string `json:"img_url"`

	Rating string `json:"rating"`
}

//Type crafted with recipe detail page in mind
type RecipeDetail struct {
	ID          string `json:"id"`
	Name        string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID  string `json:"meal_type_id"`
	ImgURL      string `json:"img_url"`
	CreatorName string `json:"creator_name"`
	Rating      string `json:"rating"`
	Steps       string `json:"steps"`
}

type RecipesHandler struct {
	DB db.DBExecutor
	SecretKey string
}

func NewRecipesHandler(db db.DBExecutor, secret string) *RecipesHandler {
	return &RecipesHandler{DB: db, SecretKey: secret}
}
