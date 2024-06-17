package main

import (
	"fmt"
	"fr/back/database"
	"fr/back/handle/fileshandlers"
	"fr/back/handle/firstpagehandler"
	"fr/back/handle/logreghandler"
	"fr/back/handle/pagehandlers"
	"fr/back/handle/posthandlers"
	"fr/back/handle/profilehandlers"
	"fr/back/handle/reactionhandlers"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	err := database.CreateDatabase()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	http.HandleFunc("/", firstpagehandler.FirstPageHandler)
	http.HandleFunc("/register", logreghandler.RegisterHandler)
	http.HandleFunc("/script.js", fileshandlers.ScriptHandler)
	http.HandleFunc("/loginForm", logreghandler.LoginFormHandler)
	http.HandleFunc("/create", posthandlers.CreateHandler)
	http.HandleFunc("/newpost", posthandlers.NewPostHandler)
	http.HandleFunc("/front/css/style.css", fileshandlers.CssHandler)
	http.HandleFunc("/login", logreghandler.LoginHadnler)
	http.HandleFunc("/sport", pagehandlers.SportPageHandler)
	http.HandleFunc("/trends", pagehandlers.TrendsPageHandler)
	http.HandleFunc("/humor", pagehandlers.HumorPageHandler)
	http.HandleFunc("/it", pagehandlers.ItPageHandler)
	http.HandleFunc("/exit", profilehandlers.ExitHandler)
	http.HandleFunc("/like", reactionhandlers.LikePostHandler)
	http.HandleFunc("/dislike", reactionhandlers.DislikePostHandler)
	http.HandleFunc("/notauthenticated", fileshandlers.CheckAuthenticatedHandler)
	http.HandleFunc("/myposts", profilehandlers.MyPostsHandler)
	http.HandleFunc("/post", posthandlers.PostPageHandler)
	http.HandleFunc("/likedposts", profilehandlers.LikedPostsHandler)
	http.HandleFunc("/addcomment", reactionhandlers.AddCommentHandler)
	http.HandleFunc("/deletepost", profilehandlers.DeletePostHandler)
	http.HandleFunc("/likecomment", reactionhandlers.LikeCommentHandler)
	http.HandleFunc("/dislikecomment", reactionhandlers.DislikeCommentHandler)
	fmt.Println("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)

}
