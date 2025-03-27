package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/dto"
)

// Errors
const (
	ErrUserNotFound = "%s: user %d not found"
)

func (s *Storage) GetUser(id int64) (*dto.User, error) {
	const op = "storage.postgres.GetUser"

	stmt, err := s.db.Prepare("SELECT id, name, email FROM public.users WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	var resultUser dto.User

	row := stmt.QueryRow(id)
	err = row.Scan(&resultUser.ID, &resultUser.Name, &resultUser.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(ErrUserNotFound, op, id)
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

	var resultUsers []*dto.User

	for rows.Next() {
		var resultUser dto.User

		err = rows.Scan(&resultUser.ID, &resultUser.Name, &resultUser.Email)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}

		resultUsers = append(resultUsers, &resultUser)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to scan rows: %w", op, err)
	}

	return resultUsers, nil
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
			log.Printf("%s: failed to close statement: %v", op, err)
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

func (s *Storage) UpdateUser(id int64, name, email string) error {
	const op = "storage.postgres.UpdateUser"

	stmt, err := s.db.Prepare("UPDATE users SET name = $1, email = $2 WHERE id = $3")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Printf("%s: failed to close statement: %v", op, err)
		}
	}()

	result, err := stmt.Exec(name, email, id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute statement: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(ErrUserNotFound, op, id)
	}

	return nil
}

func (s *Storage) DeleteUser(id int64) error {
	const op = "storage.postgres.DeleteUser"

	stmt, err := s.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Printf("%s: failed to close statement: %v", op, err)
		}
	}()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: failed to execute statement: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf(ErrUserNotFound, op, id)
	}

	return nil
}
