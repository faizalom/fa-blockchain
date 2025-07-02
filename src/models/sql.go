package models

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq" // PostgreSQL driver
)

func Conn() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("SQL_DB"))
	//db, err := sql.Open("postgres", os.Getenv("SQL_DB"))
	if err != nil {
		log.Panic(err.Error())
	}
	return db
}

// func Count(db *sql.DB, query string, args ...any) (int, error) {
// 	var count int
// 	err := db.QueryRow(query, args...).Scan(&count)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return count, err
// }

type User struct {
	ID        int64  `json:"id"`
	GoogleID  string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}

// InsertUser inserts a new user into the database
func InsertUser(user User) (int64, error) {
	db := Conn()
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func InsertUser1(user User) (int64, error) {
	db := Conn()

	query := "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}
