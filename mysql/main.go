package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func createUserTable(db *sql.DB) {
	query := `
    CREATE TABLE users (
        id INT AUTO_INCREMENT,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME,
        PRIMARY KEY (id)
    );`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func insertUser(db *sql.DB) {
	username := "third user"
	password := "this is awesome"
	createdAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
	if err != nil {
		log.Fatal(err)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(userID)

}

func queryUserById(db *sql.DB, userId int) {
	var (
		id        int
		username  string
		password  string
		createdAt time.Time
	)

	// Query the database and scan the values into out variables.
	query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
	err := db.QueryRow(query, userId).Scan(&id, &username, &password, &createdAt)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("id: %d\nusername: %s\npassword: %s\ncreatedAt: %s\n", id, username, password, createdAt)
}

func getAllUsers(db *sql.DB) {
	type user struct {
		id        int
		username  string
		password  string
		createdAt time.Time
	}

	rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("id: %d\nusername: %s\npassword: %s\ncreatedAt: %s\n\n", u.id, u.username, u.password, u.createdAt)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func deleteUserById(db *sql.DB, userId int) {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, userId)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func main() {
	loadEnvVars()

	// Configure the database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(127.0.0.1:3306)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
		return
	}

	// Initialize the first connection to the database, to see if everything works correctly.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}

	// createUserTable(db)
	// insertUser(db)
	// queryUserById(db, 2)
	// getAllUsers(db)
	deleteUserById(db, 1)
}
