package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amankp-zop/wallet/internal/api/middleware"
	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type TransactionHandler struct {
	transactionService domain.TransactionService
	validate           *validator.Validate
}

func NewTransactionHandler(ts domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: ts,
		validate:           validator.New(),
	}
}

type CreateTransferRequest struct {
	ToUserID int64           `json:"to_user_id" validate:"required"`
	Amount   decimal.Decimal `json:"amount"`
}

func (h *TransactionHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	senderUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		http.Error(w, "failed to get user id from context", http.StatusInternalServerError)
		return
	}

	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		http.Error(w, "validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}

	if senderUserID == req.ToUserID {
		http.Error(w, "Sender and receiver cannot be the same", http.StatusBadRequest)
		return
	}

	_, err := h.transactionService.CreateTransfer(r.Context(), senderUserID, req.ToUserID, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			http.Error(w, "Wallet not found for sender or receiver", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to create transfer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"Message": "success",
	})
}
