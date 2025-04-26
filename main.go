package main

import (
	"database/sql"
	"fmt"
	"go-Road/tools"
	"log"
	_ "modernc.org/sqlite"
)

const createTableSQL = `
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`

func createTable(db *sql.DB) {
	log.Println("Creating/verifying table 'clients'")
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Table 'clients' created/verified")
}

func addInitialClients(db *sql.DB, clients []tools.Client) {
	log.Println("\nAdding initial clients:")
	for _, client := range clients {
		id, err := tools.AddClient(db, client)
		if err != nil {
			log.Printf("Error adding client %s: %v\n", client.Name, err)
		} else {
			log.Printf("Added client: ID=%d, Name=%s, Email=%s\n", id, client.Name, client.Email)
		}
	}
}

func testDuplicateHandling(db *sql.DB) {
	log.Println("\nTesting duplicate email handling:")
	duplicate := tools.Client{Name: "John Second", Email: "john.doe@example.com"}
	_, err := tools.AddClient(db, duplicate)
	if err != nil {
		log.Printf("Correctly prevented duplicate: %v\n", err)
	} else {
		log.Println("ERROR: Duplicate email was incorrectly added!")
	}
}

func testClientOperations(db *sql.DB) {
	log.Println("\nClient operations completed.")

	log.Println("scanning clients by ClientID:")

	log.Println("scanning for ID:1..")
	foundClient1, err := tools.GetClientByID(db, 1)
	if err != nil {
		log.Printf("ERROR: Failed to get client by ID: %v", err)
	} else {
		log.Printf("Found client: %+v\n", foundClient1)
	}

	log.Println("\nGetting client with ID 2:")
	foundClient2, err := tools.GetClientByID(db, 2)
	if err != nil {
		log.Printf("-> Error in main getting client 2: %v\n", err)
	} else {
		log.Printf("-> Found in main: %+v\n", foundClient2)
	}

	foundClient99, err := tools.GetClientByID(db, 99)
	if err != nil {
		log.Printf("-> Error in main getting client 99: %v\n", err)
	} else {
		log.Printf("-> Found in main: %+v\n", foundClient99)
	}

	log.Println("\nGetting all clients:")
	allClients, err := tools.GetAllClientsID(db)
	if err != nil {
		log.Printf("-> Error in main getting all clients: %v\n", err)
	} else {
		log.Printf("-> Found %d clients in main:\n", len(allClients))
		for i, client := range allClients {
			log.Printf("  [%d]: %+v\n", i, client)
		}
	}

	log.Println("Testing for delete few rows.....")

	log.Println("Deleting Bob Brown (ID=4)...")
	err = tools.DeleteClient(db, 4)
	if err != nil {
		log.Printf("-> Error returned to main deleting client 4: %v\n", err)
	} else {
		log.Println("-> Delete successful for Bob Brown.")

		_, errGetBob := tools.GetClientByID(db, 4)
		if errGetBob != nil {
			log.Printf("   Verified: Bob Brown (ID=4) not found, as expected.\n")
		} else {
			log.Println("   ERROR: Bob Brown (ID=4) still found after delete!")
		}
	}

	log.Println("\nAttempting to update non-existent client (ID=99)...")
	nonExistentUpdate := tools.Client{
		ID:    99,
		Name:  "Nobody",
		Email: "nobody@nowhere.com",
	}
	err = tools.UpdateClient(db, nonExistentUpdate)
	if err != nil {
		log.Printf("-> Expected error received in main updating client 99: %v\n", err)
	} else {
		log.Println("ERROR: Updated client 99 unexpectedly!")
	}

	log.Println("\nAttempting to delete non-existent client (ID=99)...")
	err = tools.DeleteClient(db, 99)
	if err != nil {
		log.Printf("-> Expected error received in main deleting client 99: %v\n", err)
	} else {
		log.Println("ERROR: Deleted client 99 unexpectedly!")
	}
}

func main() {
	db, err := sql.Open("sqlite", "./clients.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Warning: failed to close database: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database successfully")

	createTable(db)

	// Clear existing data to start fresh
	if err := tools.ClearExistingClients(db); err != nil {
		log.Printf("Warning: %v", err)
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

	addInitialClients(db, clients)
	testDuplicateHandling(db)
	testClientOperations(db)
}
