package handler

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/service"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService domain.UserService
	validate *validator.Validate
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate: validator.New(),
	}
}

type SignupRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *UserHandler)Signup(w http.ResponseWriter, r *http.Request){
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req);err!=nil{
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err:= h.validate.Struct(req); err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err:= h.userService.Signup(r.Context(), req.Name, req.Email, req.Password)
	if err != nil{
		if errors.Is(err, service.ErrUserAlreadyExists){
			http.Error(w, err.Error(), http.StatusConflict)

			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}