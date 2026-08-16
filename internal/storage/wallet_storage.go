package storage

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type WalletStorage interface {
	TransferTx(pengirimID int, penerimaID int, jumlah int) error
}

type walletStorageImpl struct {
	db *sqlx.DB
}

func NewWalletStorage(db *sqlx.DB) WalletStorage {
	return &walletStorageImpl{db: db}
}

func (s *walletStorageImpl) TransferTx(pengirimID int, penerimaID int, jumlah int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// mencegah usr "double spend" jika user klik secar aberuntun
	var saldoPengirim int
	err = tx.Get(&saldoPengirim, "SELECT balance FROM wallets where user_id = ? FOR UPDATE", pengirimID)

	if err != nil {
		return errors.New("Error DB : " + err.Error())
	}

	if saldoPengirim < jumlah {
		return errors.New("Saldo tidak cukup")
	}
	_, err = tx.Exec("UPDATE wallets SET balance = balance - ? WHERE user_id = ?", jumlah, pengirimID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE wallets SET balance = balance + ? where user_id = ?", jumlah, penerimaID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
