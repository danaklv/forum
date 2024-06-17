package profilehandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fr/back/check"
	"net/http"
)

func ExitHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// Удаляем сессию из базы данных
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	defer database.Close()

	stmt, err := database.Prepare("DELETE FROM sessions WHERE cookie = ?")
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cookie.Value)
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusFound)

}
