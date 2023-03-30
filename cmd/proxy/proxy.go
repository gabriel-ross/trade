package main

import (
	"fmt"
	"time"

	"github.com/gabriel-ross/trade/proxy"
)

func main() {
	s := proxy.NewServer(proxy.Config{
		NAME:           "proxy",
		PORT:           "8081",
		SERVER_ADDRESS: "localhost:8080",
		CACHE_TIMEOUT:  10 * time.Second,
	})

	fmt.Printf("%v", s.Run())
}
