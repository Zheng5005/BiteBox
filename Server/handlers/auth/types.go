package auth

import "github.com/Zheng5005/BiteBox/db"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	URLPhoto string `json:"url_photo"`
	GoogleID string `json:"google_id"`
}

type AuthHandler struct {
	DB db.DBExecutor
	SecretKey string
}

func NewAuthHandler(db db.DBExecutor, secret string) *AuthHandler {
	return &AuthHandler{DB: db, SecretKey: secret}
}
