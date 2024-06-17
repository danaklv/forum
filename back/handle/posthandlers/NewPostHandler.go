package posthandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"net/http"
	"strings"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Ошибка открытия базы данных:", err)
		check.Status500(w)
		return
	}
	defer database.Close()

	title := r.FormValue("title")
	full_text := r.FormValue("full_text")
	abstract := r.FormValue("abstract")
	categories := r.Form["category"]
	category := strings.Join(categories, ",")

	if title != "" && full_text != "" && abstract != "" {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println("Ошибка получения куки:", err)
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}
		var userID int64
		err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Сессия не найдена для куки:", cookie.Value)
				http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
				return
			} else {
				fmt.Println("Ошибка выполнения запроса:", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
		}
		insertQuery := "INSERT INTO posts(title, full_text, category, likes, dislikes, user_id, abstract) VALUES (?, ?, ?, ?, ?, ?, ?)"
		statement, err := database.Prepare(insertQuery)
		if err != nil {
			fmt.Println("Ошибка подготовки запроса:", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
		defer statement.Close()
		_, err = statement.Exec(title, full_text, category, 0, 0, userID, abstract)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса:", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Название и текст поста не должны быть пустыми", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
