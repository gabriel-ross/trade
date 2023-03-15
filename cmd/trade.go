package main

import (
	"fmt"

	"github.com/gabriel-ross/trade/app"
)

func main() {
	app := app.New(app.Config{
		PORT:       "8080",
		DB_ADDRESS: "http://localhost:8529",
		DB_NAME:    "trade",
	}, app.WithCreateOnNotExist(true))

	fmt.Printf("%v", app.Run())
}
