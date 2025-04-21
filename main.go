package main

import (
	"database/sql"
	"errors"
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

func getClientByID(db *sql.DB, id int) (Client, error) {
	selectSQL := `SELECT id, name, email FROM clients WHERE id = ?`
	row := db.QueryRow(selectSQL, id)
	var client Client
	err := row.Scan(&client.ID, &client.Name, &client.Email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("INFO: No client found with ID %d\n", id)

			return Client{}, fmt.Errorf("client with ID %d not found", id)
		}

		log.Printf("ERROR: Failed to scan client data for ID %d: %v\n", id, err)
		return Client{}, err
	}
	log.Printf("SUCCESS: Found client ID=%d, Name=%s, Email=%s\n", client.ID, client.Name, client.Email)
	return client, nil

}

func clearExistingClients(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM clients")
	if err != nil {
		return fmt.Errorf("failed to clear existing clients: %w", err)
	}
	return nil
}

func getAllClientsID(db *sql.DB) ([]Client, error) {
	selectSQL := `SELECT id, name, email FROM clients`
	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all clients: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Warning: failed to close rows:", err)
		}
	}(rows)

	var clients []Client

	for rows.Next() {
		var client Client

		if err := rows.Scan(&client.ID, &client.Name, &client.Email); err != nil {

			log.Printf("ERROR: Failed to scan row while getting all clients: %v\n", err)
			return nil, err
		}

		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("ERROR: Failed to scan row while getting all clients: %v\n", err)
		return nil, err
	}
	log.Printf("SUCCESS: Found %d clients", len(clients))
	return clients, nil

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

	fmt.Println("scanning clients by ClientID:")

	fmt.Println("scanning for ID:1..")
	if foundClient1, err := getClientByID(db, 1); err != nil {
		fmt.Println("ERROR: Failed to get client by ID:", err)
	} else {
		fmt.Printf("Found client: ID=%d, Name=%s, Email=%s\n", foundClient1.ID, foundClient1.Name, foundClient1.Email)
	}

	fmt.Println("\nGetting client with ID 2:")
	foundClient2, err := getClientByID(db, 2)
	if err != nil {
		log.Printf("-> Error in main getting client 2: %v\n", err)
	} else {
		fmt.Printf("-> Found in main: %+v\n", foundClient2)
	}

	foundClient99, err := getClientByID(db, 99)
	if err != nil {
		log.Printf("-> Error in main getting client 99: %v\n", err)
	} else {
		fmt.Printf("-> Found in main: %+v\n", foundClient99)
	}
}
