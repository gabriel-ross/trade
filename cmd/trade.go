package main

import "github.com/gabriel-ross/trade"

func main() {
	app := trade.New(trade.Config{
		PORT:       "8080",
		DB_ADDRESS: "http://localhost:8529",
		DB_NAME:    "trade",
	}, trade.WithCreateOnNotExist(true))

	app.Run()
}
