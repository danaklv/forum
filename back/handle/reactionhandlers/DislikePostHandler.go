package reactionhandlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fr/back/check"
	"net/http"
	"time"
)

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		check.Status405(w)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("reiurheihtie")
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

	postID := r.FormValue("post_id")

	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Post ID is required"})
		return
	}

	var userID int64
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println(err)

	}
	var like_count int
	err = database.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&like_count)
	if err != nil {

		fmt.Println(err)
	}
	var dislike_count int
	err = database.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&dislike_count)
	if err != nil {
		fmt.Println(err)
	}
	var react string
	_ = database.QueryRow("SELECT reaction FROM user_reactions WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&react)
	defer database.Close()

	if react == "like" {
		_, err = database.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("UPDATE user_reactions SET reaction = 'dislike' WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
	} else if react == "dislike" {
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("DELETE FROM user_reactions WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			fmt.Println("hiccccc")
		}
	} else {
		_, err = database.Exec("INSERT INTO user_reactions (user_id, post_id, reaction) VALUES (?, ?, ?)", userID, postID, "dislike")
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var dislikes int
	err = database.QueryRow("SELECT dislikes FROM posts WHERE id = ?", postID).Scan(&dislikes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var likes int
	err = database.QueryRow("SELECT likes FROM posts WHERE id = ?", postID).Scan(&likes)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})

}
