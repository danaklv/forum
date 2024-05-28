package models

type User struct {
	id       int64
	username string
	email    string
	password string
}

type Post struct {
	Id       int64
	Title    string
	FullText string
	Category string
	Likes    int64
	Dislikes int64
	UserId   int64
	Abstract string
}

type Like struct {
	Id     int64
	PostId int64
	UserId int64
}

type Comment struct {
	Id       int
	PostId   int
	UserId   int
	Username string
	Text     string
}
