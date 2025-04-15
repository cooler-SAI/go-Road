package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Client struct {
	ID    int
	Name  string
	Email string
}

func main() {
	db, err := sql.Open("sqlite", "./clients.db")
	if err != nil {
		log.Fatal("We have an Error here:", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Close db Error:", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connection to DB (Ping):", err)
	}

	fmt.Println("Connected to DB")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`

	_, err = db.Exec(createTableSQL)

	if err != nil {

		log.Fatal("Error created table:", err)
	}

	fmt.Println("Table 'clients' created successful (or already created).")

}
