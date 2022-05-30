package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service{
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error){
	// struct dari user
	user := User{}
	// melakukan hash password
	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil{
		return user, err
	}

	// mapping input dari user(strunct RegisterUser) ke dalam struct User
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	user.PasswordHash = string(PasswordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil{
		return newUser, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginUserInput) (User, error){
	// ambil email dan passeeord dari user
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil{
		return user, err
	}

	// cek email
	if user.ID == 0 {
		return user, errors.New("Email not register")
	}

	// cek password
	result := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if result != nil {
		return user, errors.New("Password no match")
	}

	return user, nil
}