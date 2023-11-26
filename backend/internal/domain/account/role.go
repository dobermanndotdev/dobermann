package account

import "fmt"

var (
	UserRoleOwner  = Role{name: "owner"}
	UserRoleAdmin  = Role{name: "admin"}
	UserRoleWriter = Role{name: "writer"}
)

type Role struct {
	name string
}

func NewRole(role string) (Role, error) {
	switch role {
	case UserRoleOwner.name:
		return UserRoleOwner, nil
	case UserRoleAdmin.name:
		return UserRoleAdmin, nil
	case UserRoleWriter.name:
		return UserRoleWriter, nil
	default:
		return Role{}, fmt.Errorf("%s is not a valid user role", role)
	}
}
