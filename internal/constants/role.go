package constants

type Role string

const (
	Admin   Role = "administrator"
	Manager Role = "manager"
	Staff   Role = "staff"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsAdmin() bool {
	return r == Admin
}
