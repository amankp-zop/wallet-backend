package service

import (
	"context"
	"errors"

	"github.com/amankp-zop/wallet/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = errors.New("user with given email already exists")

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Signup(ctx context.Context, name, email, password string) (*domain.User,error){
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err!=nil{
		return nil,err
	}

	if existingUser != nil{
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return nil, err
	}

	newUser := &domain.User{
		Name: name,
		Email: email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil{
		return nil, err
	}

	return newUser, nil
}