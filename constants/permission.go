package constants

type Permission string

const (
	PermissionCreateUser     Permission = "create_user"
	PermissionReadUser       Permission = "read_user"
	PermissionUpdateUser     Permission = "update_user"
	PermissionCreateApproval Permission = "create_approval"
	PermissionReadApproval   Permission = "read_approval"
	PermissionActionApproval Permission = "action_approval"
)
