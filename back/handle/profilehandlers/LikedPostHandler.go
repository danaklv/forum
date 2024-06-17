package profilehandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"fr/back/models"
	"html/template"
	"net/http"
	"sort"
	"time"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
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
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		check.Status500(w)
		fmt.Println(err)
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
	tmpl, err := template.ParseFiles("front/templates/likedposts.html", "front/templates/header.html", "front/templates/topheader.html")
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}

	res, err := database.Query("SELECT * FROM `posts`")

	if err != nil {
		fmt.Println(err)
		check.Status500(w)
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
			check.Status401(w)
			return
		} else {
			check.Status500(w)
			return
		}
	}
	for res.Next() {
		var post models.Post
		var react string
		_ = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
		_ = database.QueryRow("SELECT reaction FROM user_reactions WHERE user_id = ? AND post_id = ?", userID, post.Id).Scan(&react)

		if react == "like" {
			posts = append(posts, post)
		}

	}
	sort.Slice(posts, func(i, j int) bool {

		return i > j
	})

	tmpl.ExecuteTemplate(w, "likedposts", posts)

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
