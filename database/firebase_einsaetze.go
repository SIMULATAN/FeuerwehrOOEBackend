package database

import (
	"context"
	"feuerwehr-ooe-backend/model"
	"log"
)

var ctx = context.Background()

func UpdateHistory(historyMap map[string]interface{}) {
	err := historyRef.Update(context.Background(), historyMap)
	if err != nil {
		log.Println("Error setting history value:", err, "\n\tPayload:", historyMap)
	}
}

func GetNotificationsToSend(einsaetze []model.Einsatz) map[string]model.Einsatz {
	var notificationsToSend = make(map[string]model.Einsatz)
	for _, einsatz := range einsaetze {
		notificationsToSend[einsatz.ID] = einsatz
	}

	// delete old notifications of finished incidents
	var existingNotificationMap map[string]bool
	err := notificationRef.OrderByKey().Get(ctx, &existingNotificationMap)
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

	return notificationsToSend
}

func UpdateSentNotifications(sentNotifications map[string]any) {
	err := notificationRef.Update(ctx, sentNotifications)
	if err != nil {
		log.Println("Error setting notification value:", err, "\n\tPayload:", sentNotifications)
	}
}
