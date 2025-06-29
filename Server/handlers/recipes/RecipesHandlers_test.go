package recipes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetRecipes_Success(t *testing.T)  {
	// Setting up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	// expected rows
	rows := sqlmock.NewRows([]string{"id", "name_recipe", "description", "meal_type_id", "img_url", "rating"}).
		AddRow("1", "Carbonara", "Best pasta in Italy", "2", "", "5").
		AddRow("2", "Pupusas", "La mejor comida de El Salvador", "1", "", "5")

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
		GROUP BY r.id
	`)).WillReturnRows(rows)

	handler := NewRecipesHandler(db, "other_key")

	req := httptest.NewRequest(http.MethodGet, "/api/recipes", nil)
	rr := httptest.NewRecorder()

	handler.RecipeHandler(rr, req)

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

func TestGetRecipe_Success(t *testing.T)  {
	// Setting up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	// expected rows
	rows := sqlmock.NewRows([]string{"id", "name_recipe", "description", "meal_type_id", "img_url", "creator_name", "avg_rating", "steps"}).
		AddRow("1", "Carbonara", "Best pasta in Italy", "2", "", "Tizio Acaso",  "5", "Put pancetta in pasta")

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT 
			r.id,
			r.name_recipe,
			r.description,
			r.meal_type_id,
			COALESCE(r.img_url, ''),
			COALESCE(u.name, r.guest_name) AS creator_name,
			COALESCE(ROUND(CAST(AVG(c.rating) AS numeric), 2), 0) AS avg_rating,
			r.steps
		FROM recipes r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN comments c ON c.recipe_id = r.id
		WHERE r.id = $1
		GROUP BY r.id, u.name, r.guest_name;
	`)).WithArgs("1").WillReturnRows(rows)

	handler := NewRecipesHandler(db, "other_key")

	req := httptest.NewRequest(http.MethodGet, "/api/recipes/1", nil)
	rr := httptest.NewRecorder()

	handler.RecipeONEHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var got RecipeDetail
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("Error decoding response %v", err)
	}

	if got.Name != "Carbonara" || got.CreatorName != "Tizio Acaso" {
		t.Errorf("Unexpected content in response: %v", got)
	}
}


// ToDo: Making a test for a request without a token = Post with auth
// ToDo: Making a test for a request wit a bad token = Post with auth 
// ToDo: Making a test for a request with a bad method = All?
// ToDo: Making a test for a request without an id = One recipe, Patch?, Delete?
// ToDo: Making a test for a request with a bad body = Post
