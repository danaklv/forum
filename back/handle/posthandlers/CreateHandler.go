package posthandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"fr/back/models"
	"html/template"
	"net/http"
	"time"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer database.Close()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			check.Status401(w)
			return
		}
		fmt.Println("Ошибка при получении cookie:", err)
		check.Status500(w)
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

	var username string
	var userID int64
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err == nil {
		err = database.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			fmt.Println(err)
		   check.Status500(w)
		}
	} else {
		fmt.Println(err)
		check.Status500(w)
	}

	tmpl, err := template.ParseFiles("front/templates/create.html", "front/templates/topheader.html", "front/templates/header.html")
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}

	tmpl.ExecuteTemplate(w, "create", nil)

	pageData := models.PageData{
		Username: username,
	}

	err = tmpl.ExecuteTemplate(w, "header", pageData)
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
	}

}
