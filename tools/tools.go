package tools

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Client struct {
	ID    int
	Name  string
	Email string
}

func AddClient(db *sql.DB, client Client) (int64, error) {
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

func GetClientByID(db *sql.DB, id int) (Client, error) {
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

func ClearExistingClients(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM clients")
	if err != nil {
		return fmt.Errorf("failed to clear existing clients: %w", err)
	}
	return nil
}

func GetAllClientsID(db *sql.DB) ([]Client, error) {
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
