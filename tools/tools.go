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

type Phone struct {
	ID          int
	ClientID    int
	PhoneNumber string
	CreatedDate string
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

func AddPhone(db *sql.DB, phone Phone) (int64, error) {
	insertSQL := `INSERT INTO phones (client_id, phone_number, created_date) VALUES (?, ?, ?)`
	result, err := db.Exec(insertSQL, phone.ClientID, phone.PhoneNumber, phone.CreatedDate)
	if err != nil {
		return 0, fmt.Errorf("failed to add phone '%s' (%s): %w",
			phone.PhoneNumber, phone.CreatedDate, err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("phone '%s' added,"+
			" but failed to get LastInsertId: %w", phone.PhoneNumber, err)
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
	log.Printf("SUCCESS: Found client %+v\n", client)
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

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Warning: failed to close rows: %v", err)
		}
	}()

	var clients []Client

	for rows.Next() {
		var client Client

		if err := rows.Scan(&client.ID, &client.Name, &client.Email); err != nil {

			log.Printf("ERROR: Failed to scan row while getting all clients: %v\n", err)
			return nil, err
		}

		clients = append(clients, client)
	}

	log.Printf("SUCCESS: Found %d clients", len(clients))
	return clients, nil

}

func DeleteClient(db *sql.DB, id int) error {
	if id == 0 {
		return fmt.Errorf("cannot delete client with zero ID")
	}
	deleteSQL := `DELETE FROM clients WHERE id = ?`
	result, err := db.Exec(deleteSQL, id)
	if err != nil {
		log.Printf("ERROR (tools): Failed to execute delete for client ID %d: %v\n", id, err)
		return fmt.Errorf("failed to delete client ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Warning (tools): Failed to get RowsAffected after delete for ID %d: %v\n", id, err)
	}

	if rowsAffected == 0 {

		log.Printf("INFO (tools): No client found with ID %d to delete.\n", id)
		return fmt.Errorf("client with ID %d not found for delete", id)
	}

	log.Printf("SUCCESS (tools): Deleted client ID=%d (RowsAffected: %d)\n", id, rowsAffected)
	return nil
}

func UpdateClient(db *sql.DB, client Client) error {

	if client.ID == 0 {
		return fmt.Errorf("cannot update client with zero ID")
	}
	updateSQL := `UPDATE clients SET name = ?, email = ? WHERE id = ?`
	result, err := db.Exec(updateSQL, client.Name, client.Email, client.ID)
	if err != nil {
		log.Printf("ERROR (tools): Failed to execute update for client ID %d: %v\n", client.ID, err)
		return fmt.Errorf("failed to update client ID %d: %w", client.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {

		log.Printf("Warning (tools): Failed to get RowsAffected after update for ID %d: %v\n", client.ID, err)
	}

	if rowsAffected == 0 {
		log.Printf("INFO (tools): No client found with ID %d to update.\n", client.ID)
		return fmt.Errorf("client with ID %d not found for update", client.ID)
	}

	log.Printf("SUCCESS (tools): Updated client %+v (RowsAffected: %d)\n",
		client, rowsAffected)
	return nil
}
