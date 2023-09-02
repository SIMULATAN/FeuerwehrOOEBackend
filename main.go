package main

import (
	"feuerwehr-ooe-backend/database"
	"feuerwehr-ooe-backend/lifecycle"
	"feuerwehr-ooe-backend/notifications"
	"feuerwehr-ooe-backend/server"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var tickerTime = 60

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!", err)
	}

	if tickerTimeStr := os.Getenv("TICKER_TIME"); tickerTimeStr != "" {
		tickerTime, err = strconv.Atoi(tickerTimeStr)
		if err != nil {
			log.Fatalf("error parsing ticker time: %v", err)
		}
	}

	database.InitFirebase()

	notifications.InitializeOneSignal()

	go server.StartServer()

	lifecycle.StartFetchThread(tickerTime)
}
