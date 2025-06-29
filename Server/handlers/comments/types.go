package comments

import (
	"github.com/Zheng5005/BiteBox/db"
)

type Comment struct {
	ID   string `json:"id"`
	UserID string `json:"user_name"`
	RecipeID string `json:"recipe_id"`
	Comment string `json:"comment"`
	Rating string `json:"rating"`
}

type CommentHandler struct {
	DB db.DBExecutor
	SecretKey string
}

func NewCommentHandler(db db.DBExecutor, secret string) *CommentHandler {
	return &CommentHandler{DB: db, SecretKey: secret}
}
