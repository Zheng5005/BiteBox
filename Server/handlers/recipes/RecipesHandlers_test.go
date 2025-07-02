package recipes

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zheng5005/BiteBox/utils"
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

func TestGetBadMethod_Sucess(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	handler := NewRecipesHandler(db, "other_key")
	req := httptest.NewRequest(http.MethodPost, "/api/recipes", nil)
	rr := httptest.NewRecorder()

	handler.RecipeHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 MetodNotAlloed, got %d", rr.Code)
	}
}

func TestPostRecipeGuest_Sucess(t *testing.T)  {
	// Setup DB mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock DB: %v", err)
	}
	defer db.Close()

	handler := NewRecipesHandler(db, "other_key")

	// Prepare multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("name", "Pupusas")
	_ = writer.WriteField("description", "Best food")
	_ = writer.WriteField("steps", "Mix and cook")
	_ = writer.WriteField("meal_type_id", "1")
	_ = writer.WriteField("guest_name", "Guesty")
	writer.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO recipes (guest_name, name_recipe, description, meal_type_id, img_url, steps)")).
		WithArgs("Guesty", "Pupusas", "Best food", "1", "", "Mix and cook").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/api/recipes/guest", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.PostRecipeGuest(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", rr.Code)
	}

	if strings.TrimSpace(rr.Body.String()) != "Recipe Created" {
		t.Errorf("Expected body 'Recipe Created', got '%s'", rr.Body.String())
	}
}

func TestPostRecipeUser_Sucess(t *testing.T)  {
	// Setup DB mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock DB: %v", err)
	}
	defer db.Close()

	handler := NewRecipesHandler(db, "other_key")

	// Prepare multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("name", "Pizza")
	_ = writer.WriteField("description", "Yummy")
	_ = writer.WriteField("steps", "Bake it")
	_ = writer.WriteField("meal_type_id", "2")
	writer.Close()

	token, err := utils.GenerateMockJWT("user-id-123", "other_key")
	if err != nil {
		t.Fatalf("Failed to generate mock JWT: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/recipes/userPost", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token) // Simulated token

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO recipes (user_id, name_recipe, description, meal_type_id, img_url, steps)")).
		WithArgs("user-id-123", "Pizza", "Yummy", "2", "", "Bake it").
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	handler.PostRecipeUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", rr.Code)
	}

	if strings.TrimSpace(rr.Body.String()) != "Recipe Created" {
		t.Errorf("Expected body 'Recipe Created', got '%s'", rr.Body.String())
	}
}
