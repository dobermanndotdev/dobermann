package account

import (
	"fmt"
	"strings"

	"errors"
)

type User struct {
	id        string
	firstName string
	lastName  string
	email     string
	role      UserRole
	accountID string
}

var (
	UserRoleOwner  = UserRole{name: "owner"}
	UserRoleAdmin  = UserRole{name: "admin"}
	UserRoleWriter = UserRole{name: "writer"}
)

type UserRole struct {
	name string
}

func NewUserRole(role string) (UserRole, error) {
	switch role {
	case UserRoleOwner.name:
		return UserRoleOwner, nil
	case UserRoleAdmin.name:
		return UserRoleAdmin, nil
	case UserRoleWriter.name:
		return UserRoleWriter, nil
	default:
		return UserRole{}, fmt.Errorf("%s is not a valid user role", role)
	}
}

func NewUser(
	id string,
	firstName string,
	lastName string,
	email string,
	role string,
	accountID string,
) (*User, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.TrimSpace(email)

	userRole, err := NewUserRole(role)
	if err != nil {
		return nil, err
	}

	if email == "" {
		return nil, errors.New("email cannot be invalid")
	}

	return &User{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		role:      userRole,
		accountID: accountID,
	}, nil
}
