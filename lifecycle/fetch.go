package lifecycle

import (
	"feuerwehr-ooe-backend/database"
	"feuerwehr-ooe-backend/model"
	"feuerwehr-ooe-backend/notifications"
	"log"
	"time"
)

func StartFetchThread(tickerTime int) {
	done := make(chan bool)

	go runFetch()
	duration := time.Second * time.Duration(tickerTime)
	ticker := time.NewTicker(duration)

	log.Println("Starting ticker with duration", duration)

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println("Tick at", t)
			runFetch()
		}
	}
}

func runFetch() {
	var historyMap = make(map[string]interface{})

	einsaetze := model.ParseLaufendeEinsaetze()

	for _, einsatz := range einsaetze {
		historyMap[einsatz.ID] = einsatz
	}

	database.UpdateHistory(historyMap)

	sendNotifications(einsaetze)
}

func sendNotifications(einsaetze []model.Einsatz) {
	sentNotifications := make(map[string]any)

	notificationsToSend := database.GetNotificationsToSend(einsaetze)

	for key, value := range notificationsToSend {
		log.Println("Sending notification for", key)
		notifications.SendOneSignalNotification(value)
		sentNotifications[key] = true
	}

	// set sent notifications to avoid sending them multiple times
	if len(sentNotifications) > 0 {
		database.UpdateSentNotifications(sentNotifications)
	}
}
