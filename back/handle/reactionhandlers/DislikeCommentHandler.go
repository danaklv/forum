package reactionhandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var expiration time.Time
	err = database.QueryRow("SELECT expiration FROM sessions WHERE cookie = ?", cookie.Value).Scan(&expiration)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	if time.Now().After(expiration) {
		_, err := database.Exec("DELETE FROM sessions WHERE cookie = ?", cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad request"})
		return
	}
	commentID := r.FormValue("comment_id")

	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "commentID is required"})
		return
	}

	defer database.Close()
	var userID int64
	_ = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	var react string
	_ = database.QueryRow("SELECT reaction FROM comment_reaction WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&react)

	if react == "dislike" {
		_, _ = database.Exec("UPDATE comment SET dislikes = dislikes - 1 WHERE id = ?", commentID)
		_, _ = database.Exec("DELETE FROM comment_reaction WHERE comment_id = ? AND user_id = ?", commentID, userID)
	} else if react == "like" {
		_, _ = database.Exec("UPDATE comment SET likes = likes - 1 WHERE id = ?", commentID)
		_, err = database.Exec("UPDATE comment_reaction SET reaction = 'dislike' WHERE comment_id = ? AND user_id = ?", commentID, userID)
		_, _ = database.Exec("UPDATE comment SET dislikes = dislikes + 1 WHERE id = ?", commentID)
	} else {
		_, _ = database.Exec("INSERT INTO comment_reaction (user_id, comment_id, reaction) VALUES (?, ?, ?)", userID, commentID, "dislike")
		_, _ = database.Exec("UPDATE comment SET dislikes = dislikes + 1 WHERE id = ?", commentID)
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var likes int
	err = database.QueryRow("SELECT likes FROM comment WHERE id = ?", commentID).Scan(&likes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var dislikes int
	err = database.QueryRow("SELECT dislikes FROM comment WHERE id = ?", commentID).Scan(&dislikes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})
}
