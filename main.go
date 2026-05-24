package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/lmittmann/tint"
	_ "modernc.org/sqlite"
)

// READ ONLY SECTOR
type Item struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DateBought     string `json:"date_bought"`
	ExpirationDate string `json:"expiration_date"`
}

var db *sql.DB
var dbPath = "go-stock.db"
var err error

// END OF READ ONLY SECTOR

func main() {
	// 1. Initialize the argparse parser first
	parser := argparse.NewParser("go-stock", "Home Stock Management Made Easy")

	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Enable detailed debug logging",
	})

	// 2. Parse the arguments
	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	// 3. Determine the log level based on the flag
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}

	// 4. Set up tint logger with the chosen level
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: "2006-01-02 15:04:05", // Added the date layout here!
	})
	slog.SetDefault(slog.New(handler))

	// --- Core logic starts here ---
	slog.Info("go-stock started")
	slog.Info("Connecting to database...")

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Successfully connected to database")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date_bought TEXT,
		expiration_date TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		slog.Error("Failed to create table", "error", err)
		os.Exit(1)
	}
	slog.Info("Inventory table verified/created")

	// Time tracking calculations
	now := time.Now()
	cleanDate := now.Format("2006-01-02")
	slog.Info("Time tracking", "right_now", cleanDate)

	expiration := time.Now().AddDate(0, 0, 14)
	cleanExpirationDate := expiration.Format("2006-01-02")
	slog.Info("Time tracking", "expires_on", cleanExpirationDate)

	Item1 := Item{
		ID:             1,
		Name:           "Milk",
		DateBought:     cleanDate,
		ExpirationDate: cleanExpirationDate,
	}

	// This debug line will show up if you pass -v or --verbose
	slog.Debug("Struct verification", "item", Item1)
}