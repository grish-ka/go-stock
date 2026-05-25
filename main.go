package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"
	"net/http"
	"html/template"
	"embed"

	"github.com/akamensky/argparse"
	"github.com/lmittmann/tint"
	_ "modernc.org/sqlite"
)

var Version = "go-stock version 0.1.0-beta.7"

// DO NOT DELETE THIS COMMENT, IT IS A SPECIAL GO COMMENT THAT ALLOWS US TO EMBED FILES INTO THE BINARY
//go:embed templates/*.html
var templateFiles embed.FS

func homeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Home Page called by", "client_ip", r.RemoteAddr)
    w.Header().Set("Content-Type", "text/html") // Tells the browser to render HTML

	switch r.Method {
		case "GET":
			// DO NOTHING, JUST LOAD THE PAGE
		
		case "POST":
			switch r.FormValue("action") {
				case "add":
					itemName := r.FormValue("name")
					// Insert the item into the database
					insertSQL := "INSERT INTO inventory (name, date_bought, expiration_date) VALUES (?, ?, ?)"
					_, err = db.Exec(insertSQL, itemName, time.Now().Format("2006-01-02"), r.FormValue("expiration_date"))
					if err != nil {
						slog.Error("Failed to insert item", "item_name", itemName, "error", err)
						fmt.Fprint(w, "Error inserting item")
					}
					slog.Debug("Item inserted successfully")
				case "delete":
					itemID := r.FormValue("id")
					deleteSQL := "DELETE FROM inventory WHERE id = ?"
					_, err = db.Exec(deleteSQL, itemID)
					if err != nil {
						slog.Error("Failed to delete item", "item_id", itemID, "error", err)
						fmt.Fprint(w, "Error deleting item")
					}
					slog.Debug("Item deleted successfully")
			}
	}

	searchQuery := r.FormValue("search")
	items, err := listItems(searchQuery)

	if err != nil {
		slog.Error("Failed to list items", "error", err)
		fmt.Fprint(w, "Error retrieving items")
		return
	}

    // Load the HTML file from the embedded memory
	tmpl, err := template.ParseFS(templateFiles, "templates/home.html")
	if err != nil {
		slog.Error("Failed to parse template file", "error", err)
		return
	}

	// Send the HTML and the database items to the browser
	tmpl.Execute(w, items)
}

func kioskHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Kiosk Page called by", "client_ip", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/html") // Tells the browser to render HTML

	switch r.Method {
		case "GET":
			// DO NOTHING, JUST LOAD THE PAGE
		
		case "POST":
			switch r.FormValue("action") {
				case "add":
					itemName := r.FormValue("name")
					// Insert the item into the database
					insertSQL := "INSERT INTO inventory (name, date_bought, expiration_date) VALUES (?, ?, ?)"
					_, err = db.Exec(insertSQL, itemName, time.Now().Format("2006-01-02"), r.FormValue("expiration_date"))
					if err != nil {
						slog.Error("Failed to insert item", "item_name", itemName, "error", err)
						fmt.Fprint(w, "Error inserting item")
					}
					slog.Debug("Item inserted successfully")
				case "delete":
					itemID := r.FormValue("id")
					deleteSQL := "DELETE FROM inventory WHERE id = ?"
					_, err = db.Exec(deleteSQL, itemID)
					if err != nil {
						slog.Error("Failed to delete item", "item_id", itemID, "error", err)
						fmt.Fprint(w, "Error deleting item")
					}
					slog.Debug("Item deleted successfully")
			}
	}

	searchQuery := r.FormValue("search")
	items, err := listItems(searchQuery)

	if err != nil {
		slog.Error("Failed to list items", "error", err)
		fmt.Fprint(w, "Error retrieving items")
		return
	}

    // Load the HTML file from the embedded memory
	tmpl, err := template.ParseFS(templateFiles, "templates/kiosk.html")
	if err != nil {
		slog.Error("Failed to parse template file", "error", err)
		return}
	// Send the HTML and the database items to the browser
	tmpl.Execute(w, items)
}

func main() {
	// 1. Initialize the argparse parser first
	parser := argparse.NewParser("go-stock", "Home Stock Management Made Easy")

	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Enable detailed debug logging",
	})

	test := parser.Flag("t", "test", &argparse.Options{
		Required: false,
		Help:     "Enable tests",
	})

	version := parser.Flag("V", "version", &argparse.Options{
		Required: false,
		Help:     "Display version information and exit",
	})

	port := parser.String("p", "port", &argparse.Options{
		Required: false,
		Help:     "Port to listen on (default: 8080)",
	})

	// 2. Parse the arguments
	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *version {
		fmt.Println(Version)
		os.Exit(0)
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

	slog.Info("----------TESTS----------")
	if *test {
		slog.Info("Running tests...")
		// Placeholder for test logic
	
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

	Item2.Name = "Bread" // This line is just to show that we can modify the struct and it won't affect the original Item1

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
} else {
	slog.Warn("Tests Skipped")
}
	slog.Info("----------WEB SERVER----------")


	// Set up HTTP handlers
	slog.Info("Setting up HTTP handlers")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/kiosk", kioskHandler)

	// Start the HTTP server
	if *port == "" {
		*port = "8080"
	}
	slog.Info("Starting HTTP server on port", "port", *port)
	err = http.ListenAndServe(":"+*port, nil)
	if err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}