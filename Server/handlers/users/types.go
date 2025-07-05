package users

import "github.com/Zheng5005/BiteBox/db"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	URLPhoto string `json:"url_photo"`
	GoogleID string `json:"google_id"`
}

type RecipesMainPage struct {
	ID   string `json:"id"`
	Name string `json:"name_recipe"`
	Description string `json:"description"`
	MealTypeID string `json:"meal_type_id"`
	ImgURL string `json:"img_url"`

	Rating string `json:"rating"`
}

type UserHandler struct {
	DB db.DBExecutor
	SecretKey string
}

func NewUserHandler(db db.DBExecutor, secret string) *UserHandler {
	return &UserHandler{DB: db, SecretKey: secret}
}
