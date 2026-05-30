package recipes

import "github.com/Zheng5005/BiteBox/db"

// RecipeSummary is a lightweight recipe representation for catalog search results.
type RecipeSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name_recipe"`
	Description string  `json:"description"`
	ImgURL      string  `json:"img_url"`
	Rating      string  `json:"rating"`
	Likes       int     `json:"likes"`
	MatchScore  float64 `json:"match_score"`
}

// CatalogRepo provides ingredient-based catalog search and AI recipe persistence.
type CatalogRepo struct {
	DB db.DBExecutor
}
