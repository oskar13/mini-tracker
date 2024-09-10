package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/oskar13/mini-tracker/pkg/data"
)

var DB *sql.DB

// Revision number in sys_info table must match this
var SchemaRevision string = "0.14"

func InitDB() {
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

func CreateSchema() error {

	// Load the schema SQL file
	schemaSQL, err := loadSQLFile("db/db_scripts/initialize_empty.sql")
	if err != nil {
		log.Fatalf("Failed to load schema file: %v", err)
		return err
	}

	// Execute the schema script
	err = executeSQL(DB, schemaSQL)
	if err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
		return err
	}

	log.Println("Database schema created successfully!")

	return nil
}

// Load SQL schema file and return its content as a string
func loadSQLFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Function to execute SQL commands
func executeSQL(db *sql.DB, sqlScript string) error {
	// MySQL driver does not support executing multiple statements by default
	statements := splitSQLStatements(sqlScript)
	for _, stmt := range statements {
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %s, error: %v", stmt, err)
		}
	}
	return nil
}

// Split SQL script into individual statements (naive approach, adjust for complex cases)
func splitSQLStatements(script string) []string {
	// Split statements by semicolon
	statements := []string{}
	for _, stmt := range strings.Split(script, ";") {
		trimmed := strings.TrimSpace(stmt)
		if len(trimmed) > 0 {
			statements = append(statements, trimmed)
		}
	}
	return statements
}
