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

I installed the library next
```bash
foo@bar:~/go-stock$ go get modernc.org/sqlite
```

The initialization of the db looks like this

```go
slog.Info("Connecting to database...")

db, err = sql.Open("sqlite", dbPath)
if err != nil {
    slog.Error("Failed to connect to database", "error", err)
    os.Exit(1)
}
defer db.Close()

slog.Info("Successfully connected to database")
```

And it creates the table automatically if it doesn't exist yet:

```go
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
```

To keep debugging easy without making the code messy, it uses argparse for flags and tint to format the logs with a clean date and time layout.

Libraries to install:

```bash
foo@bar:~/go-stock$ go get github.com/akamensky/argparse github.com/lmittmann/tint
```
The flag setup and logger initialization look like this:

```go
parser := argparse.NewParser("go-stock", "Home Stock Management Made Easy")
verbose := parser.Flag("v", "verbose", &argparse.Options{
    Required: false,
    Help:     "Enable detailed debug logging",
})

err = parser.Parse(os.Args)

logLevel := slog.LevelInfo
if *verbose {
    logLevel = slog.LevelDebug
}

handler := tint.NewHandler(os.Stdout, &tint.Options{
    Level:      logLevel,
    TimeFormat: "2006-01-02 15:04:05",
})
slog.SetDefault(slog.New(handler))
```