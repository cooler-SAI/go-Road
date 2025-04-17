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

func addClient(db *sql.DB, client Client) (int64, error) {
	insertSQL := `INSERT INTO clients (name, email) VALUES (?, ?)`
	result, err := db.Exec(insertSQL, client.Name, client.Email)
	if err != nil {
		log.Printf("ERROR: Failed to add client '%s' (%s): %v\n", client.Name, client.Email, err)
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("WARNING: Client '%s' added, but failed to get LastInsertId: %v\n", client.Name, err)
		return 0, nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("WARNING: Failed to get RowsAffected for '%s': %v\n", client.Name, err)
	}

	log.Printf("SUCCESS: Added client ID=%d, Name=%s, Email=%s (RowsAffected: %d)\n", lastInsertID, client.Name, client.Email, rowsAffected)
	return lastInsertID, nil
}

func main() {
	db, err := sql.Open("sqlite", "./clients.db")
	if err != nil {
		log.Fatal("We have an Error here:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("Close db Error:", err) // log.Println instead of fmt.Println for consistency
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
		log.Fatal("Error creating table:", err)
	}
	fmt.Println("Table 'clients' created successfully (or already exists).") // Slightly improved message

	fmt.Println("\nAdding clients via function...")

	client1 := Client{Name: "John Doe", Email: "john.doe@example.com"}
	_, err = addClient(db, client1)
	if err != nil {
		log.Println("-> Error processing John Doe in main:", err)
	}

	client2 := Client{Name: "Jane Smith", Email: "jane.smith@example.com"}
	_, err = addClient(db, client2)
	if err != nil {
		log.Println("-> Error processing Jane Smith in main:", err)
	}

	fmt.Println("\nAttempting to add client with existing email...")
	duplicateClient := Client{Name: "John Second", Email: "john.doe@example.com"}
	_, err = addClient(db, duplicateClient)
	if err != nil {
		log.Println("-> Expected error received in main when adding duplicate.")
	} else {
		log.Println("!!! LOGIC ERROR: Duplicate email was added!")
	}

	fmt.Println("\nClient operations finished.")
}
