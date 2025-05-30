package main

import (
	"log"
	"net/http"
	"github.com/Zheng5005/BiteBox/db"
	"github.com/Zheng5005/BiteBox/handlers/users"
)

func main() {
	db.InitDB()

	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong!"))
	})

	http.HandleFunc("/api/users", users.GetUsers)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

