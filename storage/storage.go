package storage

import (
	"github.com/SaidovZohid/medium_user_service/storage/postgres"
	"github.com/SaidovZohid/medium_user_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
}

type StoragePg struct {
	userRepo     repo.UserStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &StoragePg{
		userRepo:     postgres.NewUser(db),
	}
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.userRepo
}