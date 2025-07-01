package comments

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Zheng5005/BiteBox/utils"
)

func TestPostComment_Success(t *testing.T) {
	// Setting up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer db.Close()

	//Expect the INSERT query
	mock.ExpectExec(`INSERT INTO comments \(user_id, recipe_id, comment, rating\)`).
		WithArgs("user-abc", "1", "Nice recipe!", 4.5).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Initializing handler with mock DB
	handler := NewCommentHandler(db, "other_key")

	// Setting up a request Body
	body := map[string]interface{}{
		"comment": "Nice recipe!",
		"rating": 4.5,
	}
	jsonBody, _ := json.Marshal(body) 
	
	req := httptest.NewRequest(http.MethodPost, "/api/comments/post/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Attaching a valid mock jwt
	token, err := utils.GenerateMockJWT("user-abc", "other_key")
	if err != nil {
		t.Fatalf("Failed to generate mock JWT: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Preparing a response recorder
	rr := httptest.NewRecorder()

	//Calling the handler
	handler.PostComment(rr, req)

	// Asserting status Code
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", rr.Code)
	}

	// Asserting response body
	expected := "Comment created"
  if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected body '%s', got '%s'", expected, rr.Body.String())
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestGetComment_Success(t *testing.T) {
	// Setting up sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	// expected rows
	rows := sqlmock.NewRows([]string{"id", "name", "recipe_id", "comment", "rating"}).AddRow("1", "Alice", "1", "Great recipe!", "5").AddRow("2", "Bob", "1", "Too spicy!", "4")

	mock.ExpectQuery("SELECT c.id, u.name, c.recipe_id, c.comment, c.rating").
		WithArgs("1").
		WillReturnRows(rows)

	// Creating handler with mocked DB
	handler := NewCommentHandler(db, "other_key")

	// Simalating a request
	req := httptest.NewRequest(http.MethodGet, "/api/comments/1", nil)
	rr := httptest.NewRecorder()

	// Calling the handler
	handler.CommentsHandler(rr, req)

	//Checking the response
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}

	var got []Comment
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("error decoding response %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(got))
	}

	if got[0].Comment != "Great recipe!" || got[1].UserID != "Bob" {
		t.Errorf("unexpected content in response: %+v", got)
	}
}

func TestGetBadMethod_Sucess(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	handler := NewCommentHandler(db, "other_key")
	req := httptest.NewRequest(http.MethodPost, "/api/comments/1", nil)
	rr := httptest.NewRecorder()

	handler.CommentsHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 MetodNotAlloed, got %d", rr.Code)
	}
}

func TestNoToken_Success (t *testing.T) {
	// Setting up sqlmock
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer db.Close()

	// Initializing handler with mock DB
	handler := NewCommentHandler(db, "other_key")

	// Setting up a request Body
	body := map[string]interface{}{
		"comment": "Nice recipe!",
		"rating": 4.5,
	}
	jsonBody, _ := json.Marshal(body) 
	
	req := httptest.NewRequest(http.MethodPost, "/api/comments/post/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	//Preparing a response recorder
	rr := httptest.NewRecorder()

	//Calling the handler
	handler.PostComment(rr, req)

	// Asserting status Code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 201 Created, got %d", rr.Code)
	}
}

func TestBadToken_Success (t *testing.T) {
	// Setting up sqlmock
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer db.Close()

	// Initializing handler with mock DB
	handler := NewCommentHandler(db, "other_key")

	// Setting up a request Body
	body := map[string]interface{}{
		"comment": "Nice recipe!",
		"rating": 4.5,
	}
	jsonBody, _ := json.Marshal(body) 
	
	req := httptest.NewRequest(http.MethodPost, "/api/comments/post/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Attaching a invalid mock jwt
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NTE0MDczNDcsImV4cCI6MTc4Mjk0MzM0NywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSJ9.PWkBwDqKuf40en6JMHDQeDH3bUQL9e1fJJedj_xUXIo"

	req.Header.Set("Authorization", "Bearer "+token)

	//Preparing a response recorder
	rr := httptest.NewRecorder()

	//Calling the handler
	handler.PostComment(rr, req)

	// Asserting status Code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Created, got %d", rr.Code)
	}
}
