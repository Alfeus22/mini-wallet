package service

import (
	"errors"

	"github.com/Alfeus22/mini-wallet/internal/storage"
)

// kontrak interfacenya dahulu
type WalletService interface {
	TransferTx(pengirimID int, penerimaID int, jumlah int) error
}

type walletServiceImpl struct {
	storage storage.WalletStorage
}

func NewWalletService(storage storage.WalletStorage) WalletService {
	return &walletServiceImpl{storage: storage}
}

func (s *walletServiceImpl) TransferTx(pengirimID int, penerimaID int, jumlah int) error {
	if jumlah <= 0 {
		return errors.New("Jumlah transfer harus lebih dari 0")
	}
	// validasi tujuan
	if pengirimID == penerimaID {
		return errors.New("transfer yang bener")
	}

	// logika tranfer
	err := s.storage.TransferTx(pengirimID, penerimaID, jumlah)
	if err != nil {
		return err
	}
	return nil
}
