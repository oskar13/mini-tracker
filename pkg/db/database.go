package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/oskar13/mini-tracker/pkg/data"
)

var DB *sql.DB

func InitDB() {
	// Capture connection properties.

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               data.ReadPassword(os.Getenv("DB_PASSWORD_FILE")),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("DB Connected!")

	if err != nil {
		log.Fatal(err)
	}

}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
