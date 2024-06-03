package models

type User struct {
	id       int64
	username string
	email    string
	password string
}
type Post struct {
    Id           int64
    Title        string
    FullText     string
    Category     string
    Likes        int64
    Dislikes     int64
    UserId       int64
    Abstract     string
    Username     string
	CommentsCount int
}


type Like struct {
	Id     int64
	PostId int64
	UserId int64
}
type PostPageData struct {
	Post         Post
	Comments     []Comment
	CommentCount int
}

type Comment struct {
	Id int64
	Username string
	Text     string
	Likes        int64
    Dislikes     int64
}

type PageData struct {
    Username string
    Posts    []Post 
}