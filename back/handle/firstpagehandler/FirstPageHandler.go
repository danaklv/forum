package firstpagehandler

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"fr/back/models"
	"net/http"
	"sort"
	"text/template"
)

func FirstPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/templates/index.html", "front/templates/header.html", "front/templates/topheader.html")
	if err != nil {
		check.Status500(w)
		fmt.Println("Error parsing templates:", err)
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		check.Status500(w)
		fmt.Println("Error opening database:", err)
		return
	}
	defer database.Close()

	res, err := database.Query("SELECT posts.*, users.username FROM posts INNER JOIN users ON posts.user_id = users.id")
	if err != nil {
		check.Status500(w)
		fmt.Println("Error querying posts:", err)
		return
	}
	defer res.Close()

	posts := []models.Post{}
	for res.Next() {
		var post models.Post
		err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract, &post.Username)
		if err != nil {
			check.Status500(w)
			fmt.Println("Error scanning post:", err)
			return
		}
		rows, err := database.Query("SELECT users.username, comment.text FROM comment JOIN users ON comment.user_id = users.id WHERE post_id = ?", post.Id)
		if err != nil {
			check.Status500(w)
			fmt.Println("Error querying comments:", err)
			return
		}
		defer rows.Close()
		var comments []models.Comment
		for rows.Next() {
			var comment models.Comment
			err := rows.Scan(&comment.Username, &comment.Text)
			if err != nil {
				check.Status500(w)
				fmt.Println("Error scanning comment:", err)
				return
			}
			comments = append(comments, comment)
		}
		commentCount := len(comments)
		post.CommentsCount = commentCount

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return i > j
	})

	var username string
	cookie, err := r.Cookie("session_id")
	if err == nil {
		var userID int64
		err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				check.Status500(w)
				fmt.Println("Error querying session:", err)
				return
			}
		} else {
			err = database.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
			if err != nil {
				check.Status500(w)
				fmt.Println("Error querying user:", err)
				return
			}
		}
	} else {
		fmt.Println("Error retrieving cookie:", err)
	}

	pageData := models.PageData{
		Username: username,
		Posts:    posts,
	}

	err = tmpl.ExecuteTemplate(w, "index", pageData)
	if err != nil {
		check.Status500(w)
		fmt.Println("Error rendering index template:", err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "header", pageData)
	if err != nil {
		check.Status500(w)
		fmt.Println("Error rendering header template:", err)
		return
	}
}
