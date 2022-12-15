package postgres

import (
	"database/sql"
	"errors"

	"github.com/SaidovZohid/medium_user_service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type permissionRepo struct {
	db *sqlx.DB
}

func NewPermission(db *sqlx.DB) repo.PermissionStorageI {
	return &permissionRepo{
		db: db,
	}
}

func (pd *permissionRepo) CheckPermission(p *repo.Permission) (bool, error) {
	query := `
		SELECT id FROM permissions 
		WHERE user_type = $1 AND resource = $2 AND action = $3
	`

	var id int64 
	err := pd.db.QueryRow(query, p.UserType, p.Resource, p.Action).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil 
		}
		return false, nil 
	}
	return true, nil 
}