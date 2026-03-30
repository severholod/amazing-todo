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
	PhoneNumber string
}

func NewUserUninitialized(fullName string, phoneNumber string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)

}

func NewUser(id int, version int, fullName string, phoneNumber string) User {
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
	phoneNumberLen := len([]rune(u.PhoneNumber))
	if phoneNumberLen < 10 || phoneNumberLen > 15 {
		return fmt.Errorf(
			"invalid `PhoneNumber` length: %d: %w",
			phoneNumberLen,
			core_errors.ErrInvalidArgument,
		)
	}
	regex := regexp.MustCompile(`^\+[0-9]+$`)
	if !regex.MatchString(u.PhoneNumber) {
		return fmt.Errorf(
			"invalid `PhoneNumber` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if (p.FullName.Set && p.FullName.Value == nil) || (p.PhoneNumber.Set && p.PhoneNumber.Value == nil) {
		return fmt.Errorf(
			"Fields 'FullName' and 'PhoneNumber' can`t be patched to NULL :%w:",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}
	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = *patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}
	*u = tmp

	return nil
}
