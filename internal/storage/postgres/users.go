package postgres

import (
	"fmt"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/dto"
)

// Errors
const (
	ErrUserNotFound = "%s: user %d not found"
)

func (s *Storage) GetUser(id uint) (*dto.User, error) {
	const op = "storage.postgres.GetUser"

	var resultUser dto.User

	result := s.db.Take(&resultUser, id)
	if result.Error != nil {
		return nil, fmt.Errorf("%s: failed to get user: %w", op, result.Error)
	}

	return &resultUser, nil
}

func (s *Storage) GetUsersList(filter *dto.UsersListFilter) ([]*dto.User, error) {
	const op = "storage.postgres.GetUsersList"

	var resultUsers []*dto.User

	if filter == nil {
		filter = &dto.UsersListFilter{
			Page:     dto.UsersPageDefault,
			PageSize: dto.UsersPageSizeDefault,
		}
	}

	result := s.db.Find(&resultUsers).Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	if result.Error != nil {
		return nil, fmt.Errorf("%s: failed to get users list: %w", op, result.Error)
	}

	return resultUsers, nil
}

func (s *Storage) AddUser(name, email string) (uint, error) {
	const op = "storage.postgres.AddUser"

	user := &dto.User{Name: name, Email: email}

	result := s.db.Create(user)
	if result.Error != nil {
		return 0, fmt.Errorf("%s: failed to create user: %w", op, result.Error)
	}

	return user.ID, nil
}

func (s *Storage) UpdateUser(id uint, name, email string) error {
	const op = "storage.postgres.UpdateUser"

	user := &dto.User{ID: id, Name: name, Email: email}

	result := s.db.Model(&user).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("%s: failed to update user: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrUserNotFound, op, id)
	}

	return nil
}

func (s *Storage) DeleteUser(id uint) error {
	const op = "storage.postgres.DeleteUser"

	result := s.db.Delete(&dto.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("%s: failed to delete user: %w", op, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrUserNotFound, op, id)
	}

	return nil
}
