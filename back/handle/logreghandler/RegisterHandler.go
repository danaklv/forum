package logreghandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fr/back/check"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		check.Status500(w)
		return
	}
	defer database.Close()

	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed")
		check.Status405(w)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "All fields are required"})
		return
	}
	var dbUsername string
	existUsername := database.QueryRow("SELECT username FROM users WHERE username = ?", username)
	err = existUsername.Scan(&dbUsername)
	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Username already taken"})
		return
	} else if err != sql.ErrNoRows {
		check.Status500(w)
		return
	}

	existEmail := database.QueryRow("SELECT email FROM users WHERE email = ?", email)
	var dbEmail string
	err = existEmail.Scan(&dbEmail)
	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Email already taken"})
		return
	} else if err != sql.ErrNoRows {
		check.Status500(w)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		check.Status500(w)
		fmt.Println(err)
		return
	}

	stmt, err := database.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		check.Status500(w)
		fmt.Println(err)
		return
	}
	res, err := stmt.Exec(username, email, hashedPassword)
	if err != nil {
		check.Status500(w)
		fmt.Println(err)
		return
	}

	userID, err := res.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user ID"})
		fmt.Println(err)
		return
	}

	sessionID, err := uuid.NewV4()
	expiration := time.Now().Add(24 * time.Hour)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// Проверка существующей сессии
	var existingSessionID string
	err = database.QueryRow("SELECT cookie FROM sessions WHERE user_id = ?", userID).Scan(&existingSessionID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	if existingSessionID == "" {

		stmt, err := database.Prepare("INSERT INTO sessions(user_id, cookie, expiration) VALUES(?, ?, ?)")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
		_, err = stmt.Exec(userID, sessionID.String(), expiration)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
	} else {
		stmt, err := database.Prepare("UPDATE sessions SET cookie = ? WHERE user_id = ?")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
		_, err = stmt.Exec(sessionID.String(), userID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
	}
	// Установка cookie в ответе
	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID.String(),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	json.NewEncoder(w).Encode(map[string]string{"success": "User registered successfully"})
}
