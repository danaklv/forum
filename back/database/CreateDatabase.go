package database

import (
	"database/sql"
	"fmt"
)

func CreateDatabase() error {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tableCreationQueries := []string{
		`CREATE TABLE IF NOT EXISTS comment_reaction (
            id INTEGER PRIMARY KEY, 
            user_id INTEGER, 
            comment_id INTEGER, 
            reaction TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS comment (
            id INTEGER PRIMARY KEY, 
            post_id INTEGER, 
            user_id INTEGER, 
            text TEXT, 
            likes INTEGER, 
            dislikes INTEGER
        )`,
		`CREATE TABLE IF NOT EXISTS dislikes (
            id INTEGER PRIMARY KEY, 
            user_id INTEGER, 
            post_id INTEGER
        )`,
		`CREATE TABLE IF NOT EXISTS likes (
            id INTEGER PRIMARY KEY, 
            user_id INTEGER, 
            post_id INTEGER
        )`,
		`CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY, 
            title TEXT, 
            full_text TEXT, 
            category TEXT, 
            likes INTEGER, 
            dislikes INTEGER, 
            user_id INTEGER, 
            abstract TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS user_reactions (
            id INTEGER PRIMARY KEY, 
            user_id INTEGER, 
            post_id INTEGER, 
            reaction TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY, 
            username TEXT, 
            email TEXT, 
            password TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY, 
            user_id INTEGER, 
            cookie TEXT, 
            expiration DATETIME
        )`,
	}

	for _, query := range tableCreationQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	return nil
}
