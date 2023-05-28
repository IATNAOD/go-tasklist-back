package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	Hash  string `json:"-" bson:"hash"`
}

type CreateUserDTO struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

type LoginUserDTO struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

func (dto *CreateUserDTO) BuildUser() (u *User, err error) {
	passwordBytes := []byte(dto.Password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return u, err
	}

	return &User{
		ID:    "",
		Email: dto.Email,
		Hash:  string(hashedPasswordBytes),
	}, nil
}

func (user *User) CompareHashAndPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))

	return err == nil
}
