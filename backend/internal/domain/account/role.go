package account

import "fmt"

var (
	RoleOwner  = Role{value: "owner"}
	RoleAdmin  = Role{value: "admin"}
	RoleWriter = Role{value: "writer"}
)

type Role struct {
	value string
}

func NewRole(role string) (Role, error) {
	switch role {
	case RoleOwner.value:
		return RoleOwner, nil
	case RoleAdmin.value:
		return RoleAdmin, nil
	case RoleWriter.value:
		return RoleWriter, nil
	default:
		return Role{}, fmt.Errorf("the role %s is not a valid role", role)
	}
}

func (r Role) IsEmpty() bool {
	return r.value == ""
}
