package model

type Wallet struct {
	ID      int `db:"id" json:"id"`
	UserID  int `db:"user_id" json:"user_id"`
	Balance int `db:"balance" json:"-"`
}
