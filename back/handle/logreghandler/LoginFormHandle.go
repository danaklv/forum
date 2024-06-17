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

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	defer database.Close()

	email := r.FormValue("email")
	password := r.FormValue("password")

	var dbPassword []byte
	var userID int64
	err = database.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &dbPassword)
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(map[string]string{"error": "No such user"})
		return
	} else if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	err = bcrypt.CompareHashAndPassword(dbPassword, []byte(password))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Wrong password"})
		return
	}

	// Проверяем существующую сессию
	var existingSessionID string
	err = database.QueryRow("SELECT cookie FROM sessions WHERE user_id = ?", userID).Scan(&existingSessionID)
	if err != nil && err != sql.ErrNoRows {
		check.Status500(w)
		return
	}

	if existingSessionID != "" {
		// Если сессия существует, удаляем ее
		stmt, err := database.Prepare("DELETE FROM sessions WHERE user_id = ?")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
		_, err = stmt.Exec(userID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			return
		}
	}

	// Создаем новую сессию
	sessionID, err := uuid.NewV4()
	expiration := time.Now().Add(24 * time.Hour)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

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

	cookie := http.Cookie{
		Name:  "session_id",
		Value: sessionID.String(),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	json.NewEncoder(w).Encode(map[string]string{"success": "login successful"})

}
