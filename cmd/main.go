package main

import (
	"fmt"

	"github.com/gabriel-ross/trade/app"
)

var (
	// ARANGODB_ADDRESS = "http://" + os.Getenv("ARANGODB_ADDRESS") + ":" + os.Getenv("ARANGODB_PORT")
	ARANGODB_ADDRESS = "http://localhost:8529"
	PORT             = "80"
)

func main() {
	app := app.New(app.Config{
		PORT:       PORT,
		DB_ADDRESS: ARANGODB_ADDRESS,
		DB_NAME:    "trade",
	}, app.WithCreateOnNotExist(true))

	fmt.Printf("%v", app.Run())
}
