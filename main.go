package main

import (
	"fmt"
	"time"
)

// READ ONLY SECTOR
type Item struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DateBought     string `json:"date_bought"`
	ExpirationDate string `json:"expiration_date"`
}
// END OF READ ONLY SECTOR

func main() {
	fmt.Println("go-stock is still work in progress,\nthis is it for know while I plan stuff\n(the plan is plan.md file)")
	fmt.Println("----------TESTS----------")

	// Getting the current date and time
	now := time.Now()
	// Formatting the date to a clean format (YYYY-MM-DD)
	cleanDate := now.Format("2006-01-02")
	fmt.Println("Right now:", cleanDate)

	// Creating a specific date (e.g., Expiration Date)
	expiration := time.Now().AddDate(0, 0, 14) // Adding 14 days to the current date
	cleanExpirationDate := expiration.Format("2006-01-02")
	fmt.Println("Expires on:", cleanExpirationDate)

	// Creating an instance of Item with the cleaned dates
	Item1 := Item{
		ID:             1,
		Name:           "Milk",
		DateBought:     cleanDate,
		ExpirationDate: cleanExpirationDate,
	}

	// Printing the Item struct to verify the values
	fmt.Printf("Item1: %+v\n", Item1)
}