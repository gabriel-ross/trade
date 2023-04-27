package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gabriel-ross/trade/proxy"
)

func main() {
	// p := proxy.New(proxy.Config{
	// 	NAME:           "Proxy server",
	// 	PORT:           os.Getenv("PORT"),
	// 	SERVER_ADDRESS: os.Getenv("SERVER_ADDRESS"),
	// 	CACHE_TIMEOUT:  time.Hour,
	// })
	// p.Run()

	PORT := os.Getenv("PORT")
	fmt.Printf("Server running on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, proxy.Ping())
}
