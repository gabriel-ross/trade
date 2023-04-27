package main

import (
	"os"
	"time"

	"github.com/gabriel-ross/trade/proxy"
)

func main() {
	p := proxy.New(proxy.Config{
		NAME:           "Proxy server",
		PORT:           os.Getenv("PORT"),
		SERVER_ADDRESS: os.Getenv("SERVER_ADDRESS"),
		CACHE_TIMEOUT:  time.Hour,
	})
	p.Run()
}
