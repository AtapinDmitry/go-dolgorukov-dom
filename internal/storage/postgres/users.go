package postgres

import (
	"database/sql"
	"dolgorukov-dom/internal/storage/dto"
	"errors"
	"fmt"
	"log"
)

// Errors
const (
	ErrUserNotFound = "user %s not found"
)

func (s *Storage) GetUser(id string) (*dto.User, error) {
	const op = "storage.postgres.GetUser"

	stmt, err := s.db.Prepare("SELECT id, name, email FROM users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	var resultUser dto.User

	row := stmt.QueryRow(id)
	err = row.Scan(&resultUser.ID, &resultUser.Name, &resultUser.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: "+ErrUserNotFound, op, id)
		}
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	err = stmt.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to close statement: %w", op, err)
	}

	return &resultUser, nil
}

func (s *Storage) GetUsersList() ([]*dto.User, error) {
	const op = "storage.postgres.GetUsersList"

	stmt, err := s.db.Prepare("SELECT id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Printf("%s: failed to close rows: %s", op, err)
		}
	}()

	if rows.Err() != nil {

	}

}

func (s *Storage) AddUser(name, email string) (int64, error) {
	const op = "storage.postgres.AddUser"

	stmt, err := s.db.Prepare("INSERT INTO users(name, email) VALUES ($1, $2)")
	if err != nil {
		return 0, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Printf("%s: failed to close statement: %w", op, err)
		}
	}()

	result, err := stmt.Exec(name, email)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to execute statement: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteUser(id string) error {
	const op = "storage.postgres.DeleteUser"
}
