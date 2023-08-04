package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	mail2 "net/mail"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func (params *CreateUserParams) Validate() []string {
	var errors []string

	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen))
	}

	if len(params.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen))
	}

	if len(params.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
	}

	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Sprintf("incorrect email address, should be: example@domain.name"))
	}

	return errors
}

func isEmailValid(email string) bool {
	mail, err := mail2.ParseAddress(email)

	return err == nil && mail.Address == email
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
