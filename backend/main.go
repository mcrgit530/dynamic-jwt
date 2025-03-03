package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:V@v159Ha@tcp(127.0.0.1:3306)/auth?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) UNIQUE,
			password VARCHAR(255)
		)
	`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
}

type User struct {
	ID       int
	Username string
	Password string
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		log.Println("Error inserting user:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE tokens_%s (token VARCHAR(512), created_at TIMESTAMP)", user.Username))
	if err != nil {
		http.Error(w, "Failed to create token table", http.StatusInternalServerError)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("INSERT INTO tokens_%s (token, created_at) VALUES (?, ?)", user.Username), tokenString, time.Now())
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}


func signin(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println("Username:", user.Username)
	fmt.Println("Password:", user.Password)

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("SELECT id, username, password FROM users WHERE username = ? AND password = ?", user.Username, user.Password).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO tokens_%s (token, created_at) VALUES (?, ?)", user.Username), tokenString, time.Now())
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func tokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	rows, err := db.Query(fmt.Sprintf("SELECT token FROM tokens_%s", username))
	if err != nil {
		http.Error(w, "Failed to fetch tokens", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var token string
		rows.Scan(&token)
		w.Write([]byte(token + "\n"))
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/signin", signin).Methods("POST")
	router.HandleFunc("/fetchTokens/{username}", tokens).Methods("GET")

	handler := cors.Default().Handler(router)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", handler)
}
