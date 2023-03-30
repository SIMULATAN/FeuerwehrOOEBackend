package main

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"log"
)

var messagingClient *messaging.Client

func InitializeMessaging() {
	ctx := context.Background()

	var err error
	messagingClient, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalf("Error getting messaging client: %v\n", err)
	}
}

const topic = "alarmierung"

func SendNotification(einsatz Einsatz) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: einsatz.Einsatztyp.Name + " in " + einsatz.Adresse.Ort,
			Body:  "Alarmstufe " + einsatz.Alarmstufe.String(),
		},
		Topic: topic,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ChannelID:    "alarmierung",
				Priority:     messaging.PriorityMax,
				Color:        "#FF0000",
				DefaultSound: true,
			},
		},
	}

	response, err := messagingClient.Send(context.Background(), message)
	if err != nil {
		log.Println("Error sending message:", err)
	}
	log.Println("Successfully sent message:", response)
}
