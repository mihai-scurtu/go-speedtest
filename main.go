package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mihai-scurtu/speedtest/client"
	"github.com/mihai-scurtu/speedtest/server"
)

func main() {
	godotenv.Load()

	runtype := "client"

	if len(os.Args) > 1 {
		runtype = os.Args[1]
	}

	fmt.Printf("%v\n", os.Args)

	switch runtype {
	case "client":
		client.Run()

	case "server":
		server.Run()

	default:
		fmt.Println("Usage: speedtest <client|server>")
		os.Exit(1)
	}
}
