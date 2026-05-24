package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	"net/http"

	"github.com/akamensky/argparse"
	"github.com/lmittmann/tint"
	_ "modernc.org/sqlite"
)



func homeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Home Page called by", "client_ip", r.RemoteAddr)
	fmt.Fprintf(w, "Welcome to Go-Stock! Your home stock management made easy.\n\n")
	fmt.Fprint(w, "All Items:\n")
	items, err := listItems()
	if err != nil {
		slog.Error("Failed to list items", "error", err)
		fmt.Fprint(w, "Error retrieving items")
		return
	}
	for _, item := range items {
		fmt.Fprint(w, printItem(item))
	}
	// Placeholder for home page logic
}

func kioskHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Kiosk Page called by", "client_ip", r.RemoteAddr)
	fmt.Fprintf(w, "Welcome to Go-Stock! Your home stock management made easy.")
	fmt.Fprintf(w, "This is the kiosk page where you can quickly log items on the screen.\n\n")
	fmt.Fprint(w, "All Items:\n")
	items, err := listItems()
	if err != nil {
		slog.Error("Failed to list items", "error", err)
		fmt.Fprint(w, "Error retrieving items")
		return
	}
	for _, item := range items {
		fmt.Fprint(w, printItem(item))
	}
	// Placeholder for kiosk page logic
}

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
	slog.Info("----------TESTS----------")
	if *verbose {
		slog.Warn("Verbose logging enabled")
	}
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
	Item2 := Item{
		ID:             2,
		Name:           "Bread",
		DateBought:     cleanDate,
		ExpirationDate: cleanExpirationDate,
	}

	// This debug line will show up if you pass -v or --verbose
	slog.Debug("Struct verification", "item", printItem(Item1))
	slog.Debug("Struct verification", "item", printItem(Item2))


	// // Insert the item into the database
	// insertSQL := "INSERT INTO inventory (name, date_bought, expiration_date) VALUES (?, ?, ?)"
	// _, err = db.Exec(insertSQL, Item1.Name, Item1.DateBought, Item1.ExpirationDate)
	// if err != nil {
	// 	slog.Error("Failed to insert item", "item_id", Item1.ID, "error", err)
	// 	os.Exit(1)
	// }
	// slog.Info("Item inserted successfully")

	// // Insert the item into the database
	// _, err = db.Exec(insertSQL, Item2.Name, Item2.DateBought, Item2.ExpirationDate)
	// if err != nil {
	// 	slog.Error("Failed to insert item", "item_id", Item2.ID, "error", err)
	// 	os.Exit(1)
	// }
	// slog.Info("Item inserted successfully")

	ExportId := 1
	exportSQL := "SELECT id, name, date_bought, expiration_date FROM inventory WHERE id = ?"
	row := db.QueryRow(exportSQL, ExportId)
	var Exported Item
	err = row.Scan(&Exported.ID, &Exported.Name, &Exported.DateBought, &Exported.ExpirationDate)
	if err != nil {
		slog.Error("Failed to scan item", "item_id", ExportId, "error", err)
		os.Exit(1)
	}

	slog.Debug("Struct Verification", "exported", printItem(Exported))

	// Compare the original Item1 with the exported item
	if Item1 == Exported {
		slog.Info("Item1 and Exported are equal and the item was extracted successfully")
	} else {
		slog.Error("Item1 and Exported differ and the item was not extracted successfully", "original", printItem(Item1), "exported", printItem(Exported))
	}

	slog.Info("----------WEB SERVER----------")


	// Set up HTTP handlers
	slog.Info("Setting up HTTP handlers")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/kiosk", kioskHandler)

	// Start the HTTP server
	slog.Info("Starting HTTP server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}