package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Alfeus22/mini-wallet/internal/config"
	"github.com/Alfeus22/mini-wallet/internal/controller"
	"github.com/Alfeus22/mini-wallet/internal/service"
	"github.com/Alfeus22/mini-wallet/internal/storage"
)

func MockAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//menyisipkan user_id = 1 ke dalam context request
		ctx := context.WithValue(r.Context(), "user_id", 2)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func main() {

	log.Print("Memulai Server E-Wallet Mini")

	db := config.InitDD()
	defer db.Close()

	WalletStorage := storage.NewWalletStorage(db)
	WalletService := service.NewWalletService(WalletStorage)
	walletController := controller.NewWalletController(WalletService)

	mux := http.NewServeMux()

	mux.HandleFunc("/transfer", MockAuthMiddleware(walletController.TransferSaldo))

	log.Print("Server berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server mati : %v", err)
	}
}
