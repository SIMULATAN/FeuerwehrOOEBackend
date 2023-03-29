package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
	"log"
	"os"
	"strconv"
	"time"
)

var app *firebase.App

var tickerTime = 60

func main() {
	ctx := context.Background()
	var err error

	if tickerTimeStr := os.Getenv("TICKER_TIME"); tickerTimeStr != "" {
		tickerTime, err = strconv.Atoi(tickerTimeStr)
		if err != nil {
			log.Fatalf("error parsing ticker time: %v", err)
		}
	}

	firebaseConfig := GetConfig()

	ao := map[string]interface{}{"token": firebaseConfig.Token}
	conf := &firebase.Config{
		DatabaseURL:  firebaseConfig.Url,
		AuthOverride: &ao,
	}

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err = firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	InitializeMessaging()

	dbClient, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	historyRef := dbClient.NewRef("/einsatz/history")
	notificationRef := dbClient.NewRef("/einsatz/notification")

	runFetchThread(historyRef, notificationRef)
}

func runFetchThread(historyRef *db.Ref, notificationRef *db.Ref) {
	ctx := context.Background()

	done := make(chan bool)

	go runFetch(historyRef, notificationRef, ctx)
	duration := time.Second * time.Duration(tickerTime)
	ticker := time.NewTicker(duration)

	log.Println("Starting ticker with duration", duration)

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println("Tick at", t)
			runFetch(historyRef, notificationRef, ctx)
		}
	}
}

func runFetch(historyRef *db.Ref, notificationRef *db.Ref, ctx context.Context) {
	var historyMap = make(map[string]interface{})

	einsaetze := parseEinsaetze()

	for _, einsatz := range einsaetze {
		historyMap[einsatz.ID] = einsatz
	}

	var notificationsToSend = make(map[string]Einsatz)
	for _, einsatz := range einsaetze {
		notificationsToSend[einsatz.ID] = einsatz
	}

	err := historyRef.Update(ctx, historyMap)
	if err != nil {
		log.Println("Error setting history value:", err, "\n\tPayload:", historyMap)
	}

	// delete old notifications of finished incidents
	var existingNotificationMap map[string]bool
	err = notificationRef.OrderByKey().Get(ctx, &existingNotificationMap)
	if err != nil {
		log.Println("Error getting existing notification value:", err)
	}

	// delete all notifications from notificationsToSend if the incident is finished (the notifications value is true)
	for key, value := range existingNotificationMap {
		if _, ok := notificationsToSend[key]; !ok && value {
			log.Println("Deleting obsolete notification for", key)
			err = notificationRef.Child(key).Delete(ctx)
			if err != nil {
				log.Println("Error deleting existing notification value:", err)
			}
		}
	}

	// remove all notifications from notificationsToSend that are already in the database
	for key := range notificationsToSend {
		if _, ok := existingNotificationMap[key]; ok {
			delete(notificationsToSend, key)
		}
	}

	sentNotifications := make(map[string]any)

	for key, value := range notificationsToSend {
		log.Println("Sending notification for", key)
		SendNotification(value)
		sentNotifications[key] = true
	}

	// set sent notifications
	if len(sentNotifications) > 0 {
		err = notificationRef.Update(ctx, sentNotifications)
		if err != nil {
			log.Println("Error setting notification value:", err, "\n\tPayload:", sentNotifications)
		}
	}
}
