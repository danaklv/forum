package pagehandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"fr/back/models"
	"html/template"
	"net/http"
	"sort"
	"strings"
)

var posts = []models.Post{}

func SportPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/templates/sport.html", "front/templates/header.html", "front/templates/topheader.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer database.Close()

	res, err := database.Query("SELECT * FROM `posts`")

	if err != nil {
		fmt.Println(err)
		return
	}

	posts = []models.Post{}

	var username string
	cookie, err := r.Cookie("session_id")
	if err == nil {
		var userID int64
		err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
		if err == nil {
			err = database.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Ошибка при получении cookie:", err)
	}

	for res.Next() {
		var post models.Post
		err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
		if err != nil {
			fmt.Println("Error fetching username:", err)
			return
		}
		var username string
		err = database.QueryRow("SELECT username FROM users WHERE id = ?", post.UserId).Scan(&username)
		if err != nil {
			fmt.Println("Error fetching username:", err)
			return
		}
		post.Username = username
		if strings.Contains(post.Category, "Sport") {
			posts = append(posts, post)
		}

	}
	sort.Slice(posts, func(i, j int) bool {

		return i > j
	})

	err = tmpl.ExecuteTemplate(w, "sport", posts)

	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}

	pageData := models.PageData{
		Username: username,
		Posts:    posts,
	}

	err = tmpl.ExecuteTemplate(w, "header", pageData)
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}

}
