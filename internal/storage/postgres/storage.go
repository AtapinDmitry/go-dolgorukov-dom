package postgres

import (
	"database/sql"
	"fmt"

	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func New(host, port, user, password, dbname string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect: %w", op, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get sql DB: %w", op, err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to ping: %w", op, err)
	}

	return &Storage{db: db, sqlDB: sqlDB}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	if err := s.sqlDB.Close(); err != nil {
		return fmt.Errorf("%s: failed to close postgres connection: %w", op, err)
	}

	return nil
}

func (s *Storage) Migrate() error {
	const op = "storage.postgres.Migrate"

	if err := s.db.AutoMigrate(&dto.User{}); err != nil {
		return fmt.Errorf("%s: failed to auto migrate DB: %w", op, err)
	}

	return nil
}
