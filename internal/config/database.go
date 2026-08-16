package config

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func InitDD() *sqlx.DB {
	err := godotenv.Load()
	if err != nil {
		log.Print("Peringatan file .env tidak ditemukan")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL belum disetting di .envb ")
	}

	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	log.Print("Dtabase berhasil connect")
	return db
}
