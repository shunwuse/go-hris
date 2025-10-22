package constants

type Permission string

type Permissions []Permission

const (
	PermissionCreateUser     Permission = "create_user"
	PermissionReadUser       Permission = "read_user"
	PermissionUpdateUser     Permission = "update_user"
	PermissionCreateApproval Permission = "create_approval"
	PermissionReadApproval   Permission = "read_approval"
	PermissionActionApproval Permission = "action_approval"
)

func (p Permission) String() string {
	return string(p)
}

// Contains checks if the permissions contain the specified permission.
func (p Permissions) Contains(permission Permission) bool {
	for _, p := range p {
		if p == permission {
			return true
		}
	}

	return false
}

// ContainsAll checks if the permissions contain all specified permissions.
func (p Permissions) ContainsAll(permissions Permissions) bool {
	for _, permission := range permissions {
		if !p.Contains(permission) {
			return false
		}
	}

	return true
}
