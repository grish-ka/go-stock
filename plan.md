# Plan
This is the plan for go-stock

*Note:* This is the logging website server

## What Data It Stores
- Name
- Date Bought
- Expiration Date

## User Interface
**Raspberry Pi Zero 2 W** with a **touchscreen** on the fridge (or wherever it needs to be logged).The reason why ***wifi*** so it can display the local server

*Coming Later:* a camera to try and identify the item with a popup saying if it is the item

An **app** and **local server** to select your item from the phone.
### App and Local Server features
- Logging system
- Log an entire **Ocado receipt pdf**
- *(app only)* **Notifications** to popup when something is expiring so you know what to eat that day


## How Will I make It?
### 1. The Database
It will use a `struct` to make the schematic of the **database**:
```go
type Item struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	DateBought     string `json:"date_bought"`
	ExpirationDate string `json:"expiration_date"`
}
```
To this I use some date creation magic to make some random data:
```go
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
```