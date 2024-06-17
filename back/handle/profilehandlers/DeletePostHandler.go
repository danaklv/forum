package profilehandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"net/http"
	"time"
)

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		check.Status405(w)
		return
	}
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		check.Status500(w)
		return
	}
	defer database.Close()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		check.Status401(w)
		return
	}
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
		fmt.Println(err)
		check.Status401(w)
		return
	}
	postID := r.FormValue("post_id")
	_, err = database.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
	if err != nil {
		fmt.Println("Ошибка при удалении поста:", err)
		check.Status500(w)
		return
	}
	_, err = database.Exec("DELETE FROM comment WHERE post_id = ?", postID)
	if err != nil {
		fmt.Println("Ошибка при удалении поста:", err)
		check.Status500(w)
		return
	}
	_, err = database.Exec("DELETE FROM user_reactions WHERE post_id = ?", postID)
	if err != nil {
		fmt.Println("Ошибка при удалении поста:", err)
		check.Status500(w)
		return
	}
	http.Redirect(w, r, "/myposts", http.StatusSeeOther)
}
