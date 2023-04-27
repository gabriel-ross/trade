package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	// app := app.New(app.Config{
	// 	PORT:       "8080",
	// 	DB_ADDRESS: "http://localhost:8529",
	// 	DB_NAME:    "trade",
	// }, app.WithCreateOnNotExist(true))

	// fmt.Printf("%v", app.Run())

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you've hit the main server"))
	})
	PORT := os.Getenv("PORT")
	fmt.Printf("Server running on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, r)
}
