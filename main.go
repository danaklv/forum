package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fr/models"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var posts = []models.Post{}
var cook string
var WrongPassword = false
var WrongLogin = false

func main() {

	http.HandleFunc("/", FirstPageHandler)  //
	http.HandleFunc("/register", registerHandler) //
	http.HandleFunc("/script.js", ScriptHandler) //
	http.HandleFunc("/loginForm", loginFormHandler) //
	http.HandleFunc("/create", CreateHandler)  //
	http.HandleFunc("/newpost", NewPostHandler) //
	http.HandleFunc("/style.css", CssHandler) //
	http.HandleFunc("/login", LoginHadnler) //
	http.HandleFunc("/sport", SportPageHandler) //
	http.HandleFunc("/trends", TrendsPageHandler) //
	http.HandleFunc("/humor", HumorPageHandler) //
	http.HandleFunc("/exit", ExitHandler) //
	http.HandleFunc("/like", LikePostHandler) // 
	http.HandleFunc("/dislike", DislikePostHandler) //
	http.HandleFunc("/notauthenticated", CheckAuthenticatedHandler) //
	http.HandleFunc("/myposts", MyPostsHandler) //
	http.HandleFunc("/post", PostPageHandler) //
	http.HandleFunc("/likedposts", LikedPostsHandler) //
	http.HandleFunc("/addcomment", AddCommentHandler)
	http.HandleFunc("/deletepost", DeletePostHandler) //
	http.HandleFunc("/likecomment", LikeCommentHandler) //
	http.HandleFunc("/dislikecomment", DislikeCommentHandler) //

	http.HandleFunc("/it", ItPageHandler) // 

	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// _, err = db.Exec("DELETE FROM users")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }


	// _, err = db.Exec("CREATE TABLE comment_reaction (id INTEGER PRIMARY KEY, user_id INTEGER, comment_id INTEGER, reaction TEXT)")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	 

	// _, err = db.Exec("CREATE TABLE IF NOT EXIST comment (id INTEGER PRIMARY KEY, post_id INTEGER, user_id INTEGER, text TEXT, likes TEXT, dislikes TEXT)")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = db.Exec("CREATE TABLE IF NOT EXIST dislikes (id INTEGER PRIMARY KEY, user_id INTEGER,  post_id INTEGER)")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = db.Exec("CREATE TABLE IF NOT EXIST likes (id INTEGER PRIMARY KEY, user_id INTEGER,  post_id INTEGER)")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = db.Exec("CREATE TABLE IF NOT EXIST posts (id INTEGER PRIMARY KEY, title TEXT, full_text TEXT, category TEXT, likes INTEGER, dislikes INTEGER, user_id INTEGER, Abstract INTEGER")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = db.Exec("CREATE TABLE IF NOT EXIST user_reactions (id INTEGER PRIMARY KEY, user_id INTEGER,  post_id INTEGER, reaction TEXT)")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// _, err = db.Exec("CREATE TABLE IF NOT EXIST users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }





}

func loginFormHandler(w http.ResponseWriter, r *http.Request) {
    // Проверяем наличие сессии с текущим user_id
    database, err := sql.Open("sqlite3", "./forum.db")
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }
    defer database.Close()

    email := r.FormValue("email")
    password := r.FormValue("password")

    var dbPassword string
    var userID int64
    err = database.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &dbPassword)
    if err == sql.ErrNoRows {
        fmt.Println("no such user")
        json.NewEncoder(w).Encode(map[string]string{"error": "No such user"})
        return
    } else if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }

    if password != dbPassword {
        // Неверный пароль
        json.NewEncoder(w).Encode(map[string]string{"error": "Wrong password"})
        return
    }

    // Проверяем существующую сессию
    var existingSessionID string
    err = database.QueryRow("SELECT cookie FROM sessions WHERE user_id = ?", userID).Scan(&existingSessionID)
    if err != nil && err != sql.ErrNoRows {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }

    if existingSessionID != "" {
        // Если сессия существует, удаляем ее
        stmt, err := database.Prepare("DELETE FROM sessions WHERE user_id = ?")
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
            return
        }
        _, err = stmt.Exec(userID)
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
            return
        }
    }

    // Создаем новую сессию
    sessionID, err := uuid.NewV4()
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }

    stmt, err := database.Prepare("INSERT INTO sessions(user_id, cookie) VALUES(?, ?)")
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }
    _, err = stmt.Exec(userID, sessionID.String())
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }

    // Установка cookie в ответе
    cookie := http.Cookie{
        Name:  "session_id",
        Value: sessionID.String(),
        Path:  "/",
    }
    http.SetCookie(w, &cookie)


	json.NewEncoder(w).Encode(map[string]string{"success": "login successful"})
	http.Redirect(w, r, "/", http.StatusFound)
}



func SportPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/sport.html", "templates/header.html", "templates/topheader.html")
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
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func TrendsPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/trends.html", "templates/header.html", "templates/topheader.html")
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
	for res.Next() {
		var post models.Post
		err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
		if err != nil {
			fmt.Println(err)
			fmt.Println("test")
			return
		}
		var username string
        err = database.QueryRow("SELECT username FROM users WHERE id = ?", post.UserId).Scan(&username)
        if err != nil {
            fmt.Println("Error fetching username:", err)
            return
        }
        post.Username = username

		if strings.Contains(post.Category, "Trends") {
			posts = append(posts, post)
		}

	}
	sort.Slice(posts, func(i, j int) bool {

		return i > j
	})

	
    err = tmpl.ExecuteTemplate(w, "trends", posts)
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func HumorPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/humor.html", "templates/header.html", "templates/topheader.html")
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
	for res.Next() {
		var post models.Post
		err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
		if err != nil {
			fmt.Println("test1")
			fmt.Println(err)
			return
		}
		var username string
        err = database.QueryRow("SELECT username FROM users WHERE id = ?", post.UserId).Scan(&username)
        if err != nil {
            fmt.Println("Error fetching username:", err)
            return
        }
        post.Username = username
		if strings.Contains(post.Category, "Humor") {
			posts = append(posts, post)
		}

	}
	sort.Slice(posts, func(i, j int) bool {

		return i > j
	})

	err = tmpl.ExecuteTemplate(w, "humor", posts)

    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Ошибка открытия базы данных:", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	title := r.FormValue("title")
	full_text := r.FormValue("full_text")
	abstract := r.FormValue("abstract")
	categories := r.Form["category"]
	category := strings.Join(categories, ",")

	if title != "" && full_text != "" && abstract != "" {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			fmt.Println("heree")
			fmt.Println("Ошибка получения куки:", err)
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}
		fmt.Println("Получена куки:", cookie.Value)

		var userID int64
		err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Сессия не найдена для куки:", cookie.Value)
				http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
				return
			} else {
				fmt.Println("Ошибка выполнения запроса:", err)
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
		}
		fmt.Println("Найден user_id:", userID)

		insertQuery := "INSERT INTO posts(title, full_text, category, likes, dislikes, user_id, abstract) VALUES (?, ?, ?, ?, ?, ?, ?)"
		statement, err := database.Prepare(insertQuery)
		if err != nil {
			fmt.Println("Ошибка подготовки запроса:", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
		defer statement.Close()

		_, err = statement.Exec(title, full_text, category, 0, 0, userID, abstract)
		if err != nil {
			fmt.Println("Ошибка выполнения запроса:", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Название и текст поста не должны быть пустыми", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {

	_, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/notauthenticated", http.StatusSeeOther)
			return
		}
		fmt.Println("Ошибка при получении cookie:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/create.html", "templates/topheader.html", "templates/header.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.ExecuteTemplate(w, "create", nil)
}

func ScriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "script.js")
}


func FirstPageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/topheader.html")
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

    res, err := database.Query("SELECT posts.*, users.username FROM posts INNER JOIN users ON posts.user_id = users.id")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Close()

    posts := []models.Post{}
    for res.Next() {
        var post models.Post
        err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract, &post.Username)
        if err != nil {
            fmt.Println(err)
            return
        }
		rows, err := database.Query("SELECT users.username, comment.text FROM comment JOIN users ON comment.user_id = users.id WHERE post_id = ?", post.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.Username, &comment.Text)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	
    pageData := models.PageData{
        Username: username,
        Posts:    posts,
    }

    err = tmpl.ExecuteTemplate(w, "index", pageData)
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	err = tmpl.ExecuteTemplate(w, "header", pageData)  // Исправлено: выполнение header.html шаблона с `pageData`
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}



func LoginHadnler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/reg.html", "templates/header.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	if WrongLogin {
		fmt.Println("wrong login")

	}
	if WrongPassword {
		fmt.Println("wrong password")
	}

	tmpl.ExecuteTemplate(w, "reg", nil)

}

func CssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/style.css")

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    database, err := sql.Open("sqlite3", "./forum.db")
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }
    defer database.Close()

    if r.Method != http.MethodPost {
        fmt.Println("Method not allowed")
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
        return
    }

    username := r.FormValue("username")
    email := r.FormValue("email")
    password := r.FormValue("password")

    if len(username) == 0 || len(email) == 0 || len(password) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "All fields are required"})
        return
    }

    existEmail := database.QueryRow("SELECT email FROM users WHERE email = ?", email)
    var dbEmail string
    err = existEmail.Scan(&dbEmail)
    if err == nil {
        fmt.Println("email already taken")
        json.NewEncoder(w).Encode(map[string]string{"error": "Email already taken"})
        return
    } else if err != sql.ErrNoRows {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
        return
    }

    stmt, err := database.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
        fmt.Println(err)
        return
    }
    res, err := stmt.Exec(username, email, password)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to register user"})
        fmt.Println(err)
        return
    }

    userID, err := res.LastInsertId()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve user ID"})
        fmt.Println(err)
        return
    }

    sessionID, err := uuid.NewV4()
    if err != nil {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }

    // Проверка существующей сессии
    var existingSessionID string
    err = database.QueryRow("SELECT cookie FROM sessions WHERE user_id = ?", userID).Scan(&existingSessionID)
    if err != nil && err != sql.ErrNoRows {
        fmt.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
        return
    }

    fmt.Println("session ----- ", existingSessionID)

    // Если сессия существует, обновляем её, иначе создаем новую
    if existingSessionID == "" {
        stmt, err := database.Prepare("INSERT INTO sessions(user_id, cookie) VALUES(?, ?)")
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
            return
        }
        _, err = stmt.Exec(userID, sessionID.String())
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
            return
        }
    } else {
        stmt, err := database.Prepare("UPDATE sessions SET cookie = ? WHERE user_id = ?")
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
            return
        }
        _, err = stmt.Exec(sessionID.String(), userID)
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
            return
        }
    }

    // Установка cookie в ответе
    cookie := http.Cookie{
        Name:  "session_id",
        Value: sessionID.String(),
        Path:  "/",
    }
    http.SetCookie(w, &cookie)

    json.NewEncoder(w).Encode(map[string]string{"success": "User registered successfully"})
}



func ExitHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")

	fmt.Println("cookie", cookie)
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	defer database.Close()

	stmt, err := database.Prepare("DELETE FROM sessions WHERE cookie = ?")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cookie.Value)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusFound)

}


func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("Unauthorized: Missing session_id cookie")
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

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Internal server error: Failed to open database")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	defer database.Close()

	var userID int64
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Unauthorized: Session not found for cookie")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		} else {
			fmt.Println("Internal server error: Failed to query user_id from sessions")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
		return
	}

	var react string
	err = database.QueryRow("SELECT reaction FROM user_reactions WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&react)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Internal server error: Failed to query reaction from user_reactions")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	fmt.Println("Reaction:", react)

	if react == "like" {
		fmt.Println("1")
		_, err = database.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("Internal server error: Failed to update likes in posts")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("DELETE FROM user_reactions WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			fmt.Println("Internal server error: Failed to delete from user_reactions")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
	} else if react == "dislike" {
		fmt.Println("2")
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("Internal server error: Failed to update dislikes in posts")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("UPDATE user_reactions SET reaction = 'like' WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			fmt.Println("Internal server error: Failed to update reaction in user_reactions")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("Internal server error: Failed to update likes in posts")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
	} else {
		fmt.Println("3")
		_, err = database.Exec("INSERT INTO user_reactions (user_id, post_id, reaction) VALUES (?, ?, ?)", userID, postID, "like")
		if err != nil {
			fmt.Println("Internal server error: Failed to insert into user_reactions")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
		_, err = database.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("Internal server error: Failed to update likes in posts")
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
			return
		}
	}

	var likes, dislikes int
	err = database.QueryRow("SELECT likes FROM posts WHERE id = ?", postID).Scan(&likes)
	if err != nil {
		fmt.Println("Internal server error: Failed to query likes from posts")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	err = database.QueryRow("SELECT dislikes FROM posts WHERE id = ?", postID).Scan(&dislikes)
	if err != nil {
		fmt.Println("Internal server error: Failed to query dislikes from posts")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
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

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("reiurheihtie")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return

	}
	var userID int64
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("reiurheihtie")
	}
	var like_count int
	err = database.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&like_count)
	if err != nil {
		
		fmt.Println(err)
		fmt.Println("reiurheihtie")
	}
	var dislike_count int
	err = database.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&dislike_count)
	if err != nil {
		fmt.Println(err)
		fmt.Println("fffff")
	}
	var react string
	_ = database.QueryRow("SELECT reaction FROM user_reactions WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&react)
	fmt.Println(react)
	defer database.Close()

	if react == "like" {
		fmt.Println("1")
		_, err = database.Exec("UPDATE posts SET likes = likes - 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("hiffffd")
		}
		_, err = database.Exec("UPDATE user_reactions SET reaction = 'dislike' WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			fmt.Println("hidddd")
		}
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("hixxxx")
		}
	} else if react == "dislike" {
		fmt.Println("2")
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("hzzzzzi")
		}
		_, err = database.Exec("DELETE FROM user_reactions WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			fmt.Println("hiccccc")
		}
	} else {
		fmt.Println("3")
		_, err = database.Exec("INSERT INTO user_reactions (user_id, post_id, reaction) VALUES (?, ?, ?)", userID, postID, "dislike")
		if err != nil {
			fmt.Println("hi")
		}
		_, err = database.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id = ?", postID)
		if err != nil {
			fmt.Println("hiii")
		}
	}

	if err != nil {
		fmt.Println("xheuhe")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var dislikes int
	err = database.QueryRow("SELECT dislikes FROM posts WHERE id = ?", postID).Scan(&dislikes)
	if err != nil {
		fmt.Println("xheuhe")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var likes int
	err = database.QueryRow("SELECT likes FROM posts WHERE id = ?", postID).Scan(&likes)
	if err != nil {
		fmt.Println("xheuhe")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"likes": likes, "dislikes": dislikes})

}


func CheckAuthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/notauthenticated.html")
}

func MyPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/notauthenticated", http.StatusSeeOther)
			return
		}
		fmt.Println("Ошибка при получении cookie:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/myposts.html", "templates/header.html", "templates/topheader.html")
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
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Сессия не найдена для куки:", cookie.Value)
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		} else {
			fmt.Println("Ошибка выполнения запроса:", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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

}


func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/notauthenticated", http.StatusSeeOther)
			return
		}
		fmt.Println("Ошибка при получении cookie:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/likedposts.html", "templates/header.html", "templates/topheader.html")
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
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {

			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		} else {

			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
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

}


func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/notauthenticated", http.StatusSeeOther)
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Ошибка при открытии базы данных:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var userID int
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("post_id")
	commentText := r.FormValue("comment")

	_, err = database.Exec("INSERT INTO comment (post_id, user_id, text) VALUES (?, ?, ?)", postID, userID, commentText)
	if err != nil {
		fmt.Println("Ошибка при добавлении комментария:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postID), http.StatusSeeOther)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/notauthenticated", http.StatusSeeOther)
		return
	}

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Ошибка при открытии базы данных:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	var userID int
	err = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("post_id")
	_, err = database.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
	if err != nil {
		fmt.Println("Ошибка при удалении поста:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/myposts", http.StatusSeeOther)
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/post.html", "templates/header.html")
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
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.Id,&comment.Username, &comment.Text, &comment.Likes, &comment.Dislikes)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ItPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/it.html", "templates/header.html", "templates/topheader.html")
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
	for res.Next() {
		var post models.Post
		err = res.Scan(&post.Id, &post.Title, &post.FullText, &post.Category, &post.Likes, &post.Dislikes, &post.UserId, &post.Abstract)
		if err != nil {
			fmt.Println("test1")
			fmt.Println(err)
			return
		}
		var username string
        err = database.QueryRow("SELECT username FROM users WHERE id = ?", post.UserId).Scan(&username)
        if err != nil {
            fmt.Println("Error fetching username:", err)
            return
        }
        post.Username = username
		if strings.Contains(post.Category, "IT") {
			posts = append(posts, post)
		}

	}
	sort.Slice(posts, func(i, j int) bool {

		return i > j
	})
	

    err = tmpl.ExecuteTemplate(w, "it", posts)
    if err != nil {
        fmt.Println(err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}



func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	cookie, err := r.Cookie("session_id")
	if err != nil {
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

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	defer database.Close()
	var userID int64
	_ = database.QueryRow("SELECT user_id FROM sessions WHERE cookie = ?", cookie.Value).Scan(&userID)
	var react string
	_ = database.QueryRow("SELECT reaction FROM comment_reaction WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&react)
	fmt.Println(react)

	if react == "like" {
		fmt.Println("1")
		_, _ = database.Exec("UPDATE comment SET likes = likes - 1 WHERE id = ?", commentID)
		_, _ = database.Exec("DELETE FROM comment_reaction WHERE comment_id = ? AND user_id = ?", commentID, userID)
	} else if react == "dislike" {
		fmt.Println("2")
		_, _ = database.Exec("UPDATE comment SET dislikes = dislikes - 1 WHERE id = ?", commentID)
		_, err = database.Exec("UPDATE comment_reaction SET reaction = 'like' WHERE comment_id = ? AND user_id = ?", commentID, userID)
		_, _ = database.Exec("UPDATE comment SET likes = likes + 1 WHERE id = ?", commentID)
	} else {
		fmt.Println("3")
		_, _ = database.Exec("INSERT INTO comment_reaction (user_id, comment_id, reaction) VALUES (?, ?, ?)", userID, commentID, "like")
		_, _ = database.Exec("UPDATE comment SET likes = likes + 1 WHERE id = ?", commentID)
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

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	cookie, err := r.Cookie("session_id")
	if err != nil {
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

	database, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
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