package database

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	"log"
)

var app *firebase.App

func InitFirebase() {
	firebaseConfig := GetConfig()

	// Open our jsonFile
	ao := map[string]interface{}{"token": firebaseConfig.Token}
	conf := &firebase.Config{
		DatabaseURL:  firebaseConfig.Url,
		AuthOverride: &ao,
	}

	var err error

	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err = firebase.NewApp(context.Background(), conf, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	dbClient, err := app.Database(context.Background())
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	historyRef = dbClient.NewRef("/einsatz/history")
	notificationRef = dbClient.NewRef("/einsatz/notification")
}
