package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Zheng5005/BiteBox/lib"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func (h *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data (10MB max for uploaded file)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	//Upload image to Cloudinary
	file, fileHeader, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		// Assume you have a cloudinary uploader function
		imageURL, err = lib.UploadToCloudinary(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Error uploading image", http.StatusInternalServerError)
			return
		}
	} else if err != http.ErrMissingFile {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	//Save user to DB
	_, err = h.DB.Exec(
		"INSERT INTO users (name, email, password, url_photo) VALUES ($1, $2, $3, $4)",
		name, email, hashedPassword, imageURL,
	)

	if err != nil {
		log.Panicln(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var userID, hashedPassword, name, url_photo string
	err := h.DB.QueryRow("SELECT id, password, name, COALESCE(url_photo, '') FROM users WHERE email = $1", input.Email).Scan(&userID, &hashedPassword, &name, &url_photo)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	 token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"name": name,
		"url_photo": url_photo,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	errENV := godotenv.Load()
	if errENV != nil {
		log.Println("No .env file available")
	}

	secret := []byte(h.SecretKey)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
