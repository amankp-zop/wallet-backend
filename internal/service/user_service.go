package service

import (
	"context"
	"errors"
	"time"

	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shopspring/decimal"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserAlreadyExists = errors.New("user with given email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrUserNotFound = errors.New("user not found")

type userService struct {
	store     repository.Store
	jwtSecret string
	tokenTTL  time.Duration
}

func NewUserService(store repository.Store, jwtsecret string) domain.UserService {
	return &userService{
		store:     store,
		jwtSecret: jwtsecret,
		tokenTTL:  24 * time.Hour,
	}
}

func (s *userService) Signup(ctx context.Context, name, email, password string) (*domain.User,error){
	var user *domain.User

	err := s.store.ExecTx(ctx, func(q *repository.Queries)error{
		exisingUser,err := s.store.GetByEmail(ctx, email)
		if err!=nil{
			return err
		}

		if exisingUser != nil{
			return ErrUserAlreadyExists
		}

		hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err!=nil{
			return err
		}

		user = &domain.User{
			Name: name,
			Email: email,
			Password: string(hashedPassword),
		}

		err = s.store.CreateUser(ctx, user)
		if err!=nil{
			return err
		}

		walletToCreate := &domain.Wallet{
			UserID: user.ID,
			Balance: decimal.NewFromInt(0),
			Currency: "USD",
		}
		err = s.store.CreateWallet(ctx, walletToCreate)
		if err!=nil{
			return err
		}

		return nil
	})

	if err!=nil{
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error){
	user, err := s.store.GetByEmail(ctx, email)
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
	user, err := s.store.GetByID(ctx, userID)
	if err!=nil{
		return nil, err
	}

	if user == nil{
		return nil, ErrUserNotFound
	}

	return user, nil
}