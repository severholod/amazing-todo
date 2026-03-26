package domain

import (
	"fmt"
	core_errors "github.com/severholod/amazing-todo/internal/core/errors"
	"regexp"
)

type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)

}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		ID:          id,
		Version:     version,
	}
}

func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` length: %d: %w`",
			fullNameLen,
			core_errors.ErrInvalidArgument,
		)
	}
	phoneNumberLen := len([]rune(*u.PhoneNumber))
	if phoneNumberLen < 10 || phoneNumberLen > 15 {
		return fmt.Errorf(
			"invalid `PhoneNumber` length: %d: %w",
			phoneNumberLen,
			core_errors.ErrInvalidArgument,
		)
	}
	regex := regexp.MustCompile(`^\+[0-9]+$`)
	if !regex.MatchString(*u.PhoneNumber) {
		return fmt.Errorf(
			"invalid `PhoneNumber` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
