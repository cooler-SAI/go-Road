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

	// Insert a new client
	newClient := Client{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	insertSQL := `
		INSERT INTO clients (name, email) VALUES (?, ?)
	`

	result, err := db.Exec(insertSQL, newClient.Name, newClient.Email)
	if err != nil {
		log.Fatal("Error inserting client:", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal("Error getting affected rows:", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Error getting last insert ID:", err)
	}

	fmt.Printf("Client inserted successfully. Rows affected: %d, Last Insert ID: %d\n", rowsAffected, lastInsertID)

	newClient2 := Client{
		Name:  "Jane Smith",
		Email: "jane.smith@example.com",
	}
	_, err = db.Exec(insertSQL, newClient2.Name, newClient2.Email)
	if err != nil {
		log.Fatal("Error inserting client:", err)
	}
	fmt.Println("Client inserted successfully.")
}
