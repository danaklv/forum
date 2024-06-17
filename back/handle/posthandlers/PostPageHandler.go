package posthandlers

import (
	"database/sql"
	"fmt"
	"fr/back/check"
	"fr/back/models"
	"html/template"
	"net/http"
)

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("front/templates/post.html", "front/templates/header.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var post models.Post
	err = database.QueryRow("SELECT * FROM posts WHERE Id = ?", id).Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
	if err != nil {
		if err == sql.ErrNoRows {
			check.Status404(w)
			return
		} else {
			fmt.Println(err)
			check.Status500(w)
			return
		}
	}
	var username string
	err = database.QueryRow("SELECT username FROM users WHERE id = ?", post.UserId).Scan(&username)
	if err != nil {
		fmt.Println("Error fetching username:", err)
		return
	}
	post.Username = username

	// Загрузка комментариев
	rows, err := database.Query("SELECT comment.id, users.username, comment.text, comment.likes, comment.dislikes FROM comment JOIN users ON comment.user_id = users.id WHERE post_id = ?", post.Id)
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.Id, &comment.Username, &comment.Text, &comment.Likes, &comment.Dislikes)
		if err != nil {
			fmt.Println(err)
		    check.Status500(w)
		    return
		}
		comments = append(comments, comment)
	}

	// Подсчет количества комментариев
	commentCount := len(comments)

	data := models.PostPageData{
		Post:         post,
		Comments:     comments,
		CommentCount: commentCount,
	}

	err = tmpl.ExecuteTemplate(w, "post", data)
	if err != nil {
		fmt.Println(err)
		check.Status500(w)
		return
	}
}
