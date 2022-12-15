package repo

type Permission struct {
	UserType string
	Resource string
	Action   string
}

type PermissionStorageI interface {
	CheckPermission(*Permission) (bool, error)
}

