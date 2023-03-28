package main

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"log"
)

var client *messaging.Client

func InitializeMessaging() {
	ctx := context.Background()

	var err error
	client, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
}

func SendNotification(einsatz Einsatz) {
	// The topic name can be optionally prefixed with "/topics/".
	topic := "alarmierung"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: einsatz.Einsatzart + " in " + einsatz.Adresse.Ort,
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

	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}
