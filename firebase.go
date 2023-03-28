package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type FirebaseConfig struct {
	Url   string `json:"firebase_url"`
	Token string `json:"firebase_token"`
}

func GetConfig() FirebaseConfig {
	jsonFile, err := os.Open("serviceAccountKey.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var firebaseConfig FirebaseConfig

	err = json.Unmarshal(bytes, &firebaseConfig)
	if err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	if firebaseConfig.Url == "" {
		log.Fatalf("error unmarshaling config: %v", "firebase_url is empty")
	}

	if firebaseConfig.Token == "" {
		log.Fatalf("error unmarshaling config: %v", "firebase_token is empty")
	}

	return firebaseConfig
}
