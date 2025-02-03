package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const SessionName string = "GO_SESSION" // Session cookie name
const UserIDKey string = "user_id"      // Key used in session

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, SessionName)
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		if user_id, ok := session.Values[UserIDKey].(uint); !ok || user_id == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authRouter(r *mux.Router) {
	r.HandleFunc("/Register", handleRegister).Methods("POST") // TODO: not in swagger
	r.HandleFunc("/Login", handleLogin).Methods("POST")
	r.HandleFunc("/Logoff", handleLogoff).Methods("POST")
	r.HandleFunc("/GetToken", handleGetToken).Methods("POST")
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("userEmail")
	password := r.URL.Query().Get("userPassword")
	if email == "" || password == "" {
		http.Error(w, "Expected userEmail and userPassword in query", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	user := User{
		Email:    email,
		Password: string(hashedPassword),
	}

	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, "Could not create user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("userEmail")
	password := r.URL.Query().Get("userPassword")
	if email == "" || password == "" {
		http.Error(w, "Expected userEmail and userPassword in query", http.StatusBadRequest)
		return
	}

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, err := store.Get(r, SessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID

	if err := session.Save(r, w); err != nil {
		log.Println(err)
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user":    user,
	})
}

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("userEmail")
	password := r.URL.Query().Get("userPassword")
	if email == "" || password == "" {
		http.Error(w, "Expected userEmail and userPassword in query", http.StatusBadRequest)
		return
	}

	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, err := store.Get(r, SessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values[UserIDKey] = user.ID

	if err := session.Save(r, w); err != nil {
		log.Println(err)
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(session.ID)) // TODO: not documented in swagger, plus ID is not sufficient to reconstruct session later on
}

func handleLogoff(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SessionName)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values[UserIDKey] = nil
	session.Options.MaxAge = -1 // Delete the session

	if err := session.Save(r, w); err != nil {
		log.Println(err)
		http.Error(w, "Could not save session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logged out successfully",
	})
}
