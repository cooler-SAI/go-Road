package main

import (
	"database/sql"
	"fmt"
	"go-Road/tools"
	"log"
	_ "modernc.org/sqlite"
)

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
	if err := tools.ClearExistingClients(db); err != nil {
		log.Println("Warning:", err)
	}

	// Define test clients
	var clients = []tools.Client{
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
		id, err := tools.AddClient(db, client)
		if err != nil {
			log.Printf("Error adding client %s: %v\n", client.Name, err)
		} else {
			fmt.Printf("Added client: ID=%d, Name=%s, Email=%s\n", id, client.Name, client.Email)
		}
	}

	// Test duplicate handling
	fmt.Println("\nTesting duplicate email handling:")
	duplicate := tools.Client{Name: "John Second", Email: "john.doe@example.com"}
	_, err = tools.AddClient(db, duplicate)
	if err != nil {
		fmt.Printf("Correctly prevented duplicate: %v\n", err)
	} else {
		log.Println("ERROR: Duplicate email was incorrectly added!")
	}

	fmt.Println("\nClient operations completed.")

	fmt.Println("scanning clients by ClientID:")

	fmt.Println("scanning for ID:1..")
	if foundClient1, err := tools.GetClientByID(db, 1); err != nil {
		fmt.Println("ERROR: Failed to get client by ID:", err)
	} else {
		fmt.Printf("Found client: ID=%d, Name=%s, Email=%s\n",
			foundClient1.ID, foundClient1.Name, foundClient1.Email)
	}

	fmt.Println("\nGetting client with ID 2:")
	foundClient2, err := tools.GetClientByID(db, 2)
	if err != nil {
		log.Printf("-> Error in main getting client 2: %v\n", err)
	} else {
		fmt.Printf("-> Found in main: %+v\n", foundClient2)
	}

	foundClient99, err := tools.GetClientByID(db, 99)
	if err != nil {
		log.Printf("-> Error in main getting client 99: %v\n", err)
	} else {
		fmt.Printf("-> Found in main: %+v\n", foundClient99)
	}

	fmt.Println("\nGetting all clients:")
	allClients, err := tools.GetAllClientsID(db)
	if err != nil {
		log.Printf("-> Error in main getting all clients: %v\n", err)
	} else {
		fmt.Printf("-> Found %d clients in main:\n", len(allClients))
		for i, client := range allClients {

			fmt.Printf("  [%d]: %+v\n", i, client)
		}
	}

	fmt.Println("Testing for delete few rows.....")

	fmt.Println("Deleting Bob Brown (ID=4)...")
	err = tools.DeleteClient(db, 4)
	if err != nil {
		log.Printf("-> Error returned to main deleting client 4: %v\n", err)
	} else {
		fmt.Println("-> Delete successful for Bob Brown.")

		_, errGetBob := tools.GetClientByID(db, 4)
		if errGetBob != nil {
			fmt.Printf("   Verified: Bob Brown (ID=4) not found, as expected.\n")
		} else {
			fmt.Println("   ERROR: Bob Brown (ID=4) still found after delete!")
		}
	}

	fmt.Println("\nAttempting to delete non-existent client (ID=99)...")
	err = tools.DeleteClient(db, 99)
	if err != nil {

		log.Printf("-> Expected error received in main deleting client 99: %v\n", err)
	} else {
		log.Println("ERROR: Deleted client 99 unexpectedly!")
	}

}
