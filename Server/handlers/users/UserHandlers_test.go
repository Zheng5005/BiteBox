package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zheng5005/BiteBox/utils"
)

func TestGetUserRecipes_Success(t *testing.T)  {
	// Setting up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	// expected rows
	rows := sqlmock.NewRows([]string{"id", "name_recipe", "description", "meal_type_id", "img_url", "rating",}).
		AddRow("1", "Carbonara", "Best pasta in Italy", "2", "", "5",).
		AddRow("2", "Pupusas", "La mejor comida de El Salvador", "1", "", "5",)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			r.id, 
			r.name_recipe, 
			r.description, 
			r.meal_type_id, 
			COALESCE(r.img_url, ''),
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg 
		FROM recipes r 
		LEFT JOIN comments c ON r.id = c.recipe_id
		WHERE r.user_id = $1
		GROUP BY r.id
	`)).WithArgs("5").WillReturnRows(rows)

	token, err := utils.GenerateMockJWT("5", "other_key")
	if err != nil {
		t.Fatalf("Failed to generate mock JWT: %v", err)
	}

	handler := NewUserHandler(db, "other_key")

	req := httptest.NewRequest(http.MethodGet, "/api/recipes", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Simulated token
	rr := httptest.NewRecorder()

	handler.GetRecipesAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var got []RecipesMainPage
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("Error decoding response %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Expected 2 recipes, got %v", len(got))
	}

	if got[0].Name != "Carbonara" || got[1].Description != "La mejor comida de El Salvador" {
		t.Errorf("Unexpected content in response: %v", got)
	}
}

func TestEditRecipe_Success(t *testing.T)  {
	
}

func TestDeactivateRecipe_Success(t *testing.T)  {
	
}

func TestActivateRecipe_Success(t *testing.T)  {
	
}

func TestGetRecipesByUser_Success(t *testing.T)  {
	
}

func TestGetRecipesByGuest_Success(t *testing.T)  {
	
}
