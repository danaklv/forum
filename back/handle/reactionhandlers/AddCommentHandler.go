package reactionhandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"net/http"
	"time"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		check.Status405(w)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		check.Status401(w)
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Ошибка при открытии базы данных:", err)
		check.Status500(w)
		return
	}
	defer database.Close()
	var expiration time.Time
	err = database.QueryRow("SELECT expiration FROM sessions WHERE cookie = ?", cookie.Value).Scan(&expiration)
	if err != nil {
		check.Status401(w)
		return
	}
	if time.Now().After(expiration) {
		_, err := database.Exec("DELETE FROM sessions WHERE cookie = ?", cookie.Value)
		if err != nil {
			check.Status500(w)
			return
		}
		check.Status401(w)
		return
	}

	var userID int
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		check.Status401(w)
		return
	}

	postID := r.FormValue("post_id")
	commentText := r.FormValue("comment")
	if commentText == "" {
		check.Status400(w)
		return
	}

	_, err = database.Exec("INSERT INTO comment (post_id, user_id, text) VALUES (?, ?, ?)", postID, userID, commentText)
	if err != nil {
		fmt.Println("Ошибка при добавлении комментария:", err)
		check.Status500(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postID), http.StatusSeeOther)
}
