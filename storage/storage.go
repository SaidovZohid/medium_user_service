package storage

import (
	"github.com/SaidovZohid/medium_user_service/storage/postgres"
	"github.com/SaidovZohid/medium_user_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Permission() repo.PermissionStorageI
}

type StoragePg struct {
	userRepo       repo.UserStorageI
	permissionRepo repo.PermissionStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &StoragePg{
		userRepo:       postgres.NewUser(db),
		permissionRepo: postgres.NewPermission(db),
	}
}

func (s *StoragePg) User() repo.UserStorageI {
	return s.userRepo
}

func (s *StoragePg) Permission() repo.PermissionStorageI {
	return s.permissionRepo
}
