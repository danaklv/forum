package profilehandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"fr/back/models"
	"html/template"
	"net/http"
	"sort"
)

var posts = []models.Post{}

func MyPostsHandler(w http.ResponseWriter, r *http.Request) {
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
	tmpl, err := template.ParseFiles("front/templates/myposts.html", "front/templates/header.html", "front/templates/topheader.html")
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
	var userID int64
	var username string
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err == nil {
		err = database.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Сессия не найдена для куки:", cookie.Value)
			check.Status401(w)
			return
		} else {
			fmt.Println("Ошибка выполнения запроса:", err)
			check.Status500(w)
			return
		}
	}
	for res.Next() {
		var post models.Post
		err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
		if err != nil {
			fmt.Println(err)
			return
		}
		if post.UserId == userID {
			posts = append(posts, post)
		}

	}
	sort.Slice(posts, func(i, j int) bool {

		return i > j
	})

	tmpl.ExecuteTemplate(w, "myposts", posts)

	pageData := models.PageData{
		Username: username,
		Posts:    posts,
	}

	err = tmpl.ExecuteTemplate(w, "header", pageData)
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
	}

}
