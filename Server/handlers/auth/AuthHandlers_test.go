package auth

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
	"golang.org/x/crypto/bcrypt"
)

func TestSignIn_Success(t *testing.T)  {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock DB: %v", err)
	}
	defer db.Close()

	handler := NewAuthHandler(db, "other_key")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("name", "Jane")
	_ = writer.WriteField("email", "jd@gmail.com")
	_ = writer.WriteField("password", "123")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email, password, url_photo)")).
		WithArgs("Jane", "jd@gmail.com", sqlmock.AnyArg(), "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()
	handler.SignUpHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", rr.Code)
	}

	if strings.TrimSpace(rr.Body.String()) != "User created" {
		t.Errorf("Expected body 'User created', got '%s'", rr.Body.String())
	}
}

func TestLogin_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock DB: %v", err)
	}
	defer db.Close()

	handler := NewAuthHandler(db, "other_key")
	
	email := "jd@gmail.com"
	password := "123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, password, name, COALESCE(url_photo, '') FROM users WHERE email = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "name", "url_photo"}).
			AddRow("1", string(hashed), "Test User", "photo.jpg"))

	body := `{"email": "jd@gmail.com", "password": "123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.LoginHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 Created, got %d", rr.Code)
	}

  var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := resp["token"]; !ok {
		t.Errorf("Expected JWT token in response, got: %v", resp)
	}
}
