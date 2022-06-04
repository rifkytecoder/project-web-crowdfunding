package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	// Implement mapping request
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user" //hardcode

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User email not founded")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {

	email := input.Email
	// Check email if registered in database
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	// Register, if the user is not registered
	if user.ID == 0 {
		return true, nil
	}

	return false, nil

}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// dapatkan user berdasarkan ID
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	// update attribute avatar file name
	user.AvatarFileName = fileLocation

	// simpan perubahan avatar file name
	updateUser, err := s.repository.Update(user)
	if err != nil {
		return updateUser, err
	}
	return updateUser, nil
}
