package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Alfeus22/mini-wallet/internal/service"
)

type WalletController struct {
	service service.WalletService
}

func NewWalletController(service service.WalletService) *WalletController {
	return &WalletController{service: service}
}

// format jsnon yang diharapkan user
type TransferRequest struct {
	PenerimaID int `json:"penerima_id"`
	Jumlah     int `json:"jumlah"`
}

func (c *WalletController) TransferSaldo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	// buka json dari user
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "format request tidak valid"})
		return
	}
	// middleware
	pengirimID, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "anda belum login"})
		return
	}

	// lempar ke dapur / service
	err := c.service.TransferTx(pengirimID, req.PenerimaID, req.Jumlah)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "berhasil transfer"})
}
