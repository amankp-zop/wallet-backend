package service

import (
	"context"
	"errors"
	"time"

	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = errors.New("user with given email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrUserNotFound = errors.New("user not found")

type userService struct {
	userRepo domain.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewUserService(userRepo domain.UserRepository, jwtsecret string) domain.UserService {
	return &userService{
		userRepo: userRepo,
		jwtSecret: jwtsecret,
		tokenTTL:  24 * time.Hour,
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

func (s *userService) Login(ctx context.Context, email, password string) (string, error){
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err!=nil{
		return "", err
	}

	if user == nil{
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	if err != nil{
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(s.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil{
		return "", err
	}

	return signedToken, nil
}

func (s *userService) GetProfile(ctx context.Context, userID int64) (*domain.User, error){
	user, err := s.userRepo.GetByID(ctx, userID)
	if err!=nil{
		return nil, err
	}

	if user == nil{
		return nil, ErrUserNotFound
	}

	return user, nil
}