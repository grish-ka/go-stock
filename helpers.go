package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	// "os"
	// "time"
	"strings"

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

func printItem(item Item) string {
	var output strings.Builder
	fmt.Fprintf(&output, "ID: %d\n", item.ID)
	fmt.Fprintf(&output, "Name: %s\n", item.Name)
	fmt.Fprintf(&output, "Date Bought: %s\n", item.DateBought)
	fmt.Fprintf(&output, "Expiration Date: %s\n", item.ExpirationDate)

	slog.Debug("Item details:", "item", output.String())

	return output.String()
}

func listItems() ([]Item, error) {
	rows, err := db.Query("SELECT id, name, date_bought, expiration_date FROM inventory")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.DateBought, &item.ExpirationDate)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// END OF READ ONLY SECTOR