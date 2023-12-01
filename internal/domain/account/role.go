package account

import "fmt"

var (
	RoleOwner  = Role{name: "owner"}
	RoleAdmin  = Role{name: "admin"}
	RoleWriter = Role{name: "writer"}
)

type Role struct {
	name string
}

func NewRole(role string) (Role, error) {
	switch role {
	case RoleOwner.name:
		return RoleOwner, nil
	case RoleAdmin.name:
		return RoleAdmin, nil
	case RoleWriter.name:
		return RoleWriter, nil
	default:
		return Role{}, fmt.Errorf("%s is not a valid user role", role)
	}
}

func (r Role) String() string {
	return r.name
}
