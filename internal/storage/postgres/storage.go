package postgres

import (
	"database/sql"
	"fmt"
	"net"

	// Postgres registering
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, password, dbname string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	storageInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, net.JoinHostPort(host, port), dbname)

	db, err := sql.Open("postgres", storageInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to ping: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: failed to close postgres connection: %w", op, err)
	}

	return nil
}
