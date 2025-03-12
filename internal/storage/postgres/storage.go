package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	storageInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", storageInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect: %v", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to ping: %v", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: failed to close postgres connection: %v", op, err)
	}

	return nil
}
