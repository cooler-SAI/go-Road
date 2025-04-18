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
		return 0, fmt.Errorf("failed to add client '%s' (%s): %w", client.Name, client.Email, err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("client '%s' added, but failed to get LastInsertId: %w", client.Name, err)
	}

	return lastInsertID, nil
}

func clearExistingClients(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM clients")
	if err != nil {
		return fmt.Errorf("failed to clear existing clients: %w", err)
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlite", "./clients.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Warning: failed to close database:", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database successfully")

	// Create table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table 'clients' created/verified")

	// Clear existing data to start fresh
	if err := clearExistingClients(db); err != nil {
		log.Println("Warning:", err)
	}

	// Define test clients
	clients := []Client{
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
		{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
		{Name: "Bob Brown", Email: "bob.brown@example.com"},
		{Name: "Eve Green", Email: "eve.green@example.com"},
		{Name: "Charlie White", Email: "charlie.white@example.com"},
		{Name: "Diana Black", Email: "diana.black@example.com"},
		{Name: "Frank Red", Email: "frank.red@example.com"},
	}

	fmt.Println("\nAdding initial clients:")
	for _, client := range clients {
		id, err := addClient(db, client)
		if err != nil {
			log.Printf("Error adding client %s: %v\n", client.Name, err)
		} else {
			fmt.Printf("Added client: ID=%d, Name=%s, Email=%s\n", id, client.Name, client.Email)
		}
	}

	// Test duplicate handling
	fmt.Println("\nTesting duplicate email handling:")
	duplicate := Client{Name: "John Second", Email: "john.doe@example.com"}
	_, err = addClient(db, duplicate)
	if err != nil {
		fmt.Printf("Correctly prevented duplicate: %v\n", err)
	} else {
		log.Println("ERROR: Duplicate email was incorrectly added!")
	}

	fmt.Println("\nClient operations completed.")
}
