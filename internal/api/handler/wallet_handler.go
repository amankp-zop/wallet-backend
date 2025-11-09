package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amankp-zop/wallet/internal/api/middleware"
	"github.com/amankp-zop/wallet/internal/domain"
	"github.com/amankp-zop/wallet/internal/service"
)

type WalletHandler struct {
	walletService domain.WalletService
}

func NewWalletHandler(walletService domain.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func(h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request){

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wallet, err := h.walletService.GetWalletByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(wallet)
}
