package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const jsonUrl string = "https://cf-intranet.ooelfv.at/webext2/rss/json_laufend.txt"

type Einsatz struct {
	ID          string      `json:"num1"`
	Einsatzort  string      `json:"einsatzort"`
	Koordinaten Koordinaten `json:"wgs84"`
	Alarmstufe  json.Number `json:"alarmstufe"`
	Startzeit   MyTime      `json:"startzeit"`
	Einsatzname string      `json:"einsatzname"`
	Einsatzart  string      `json:"einsatzart"`
	Einsatztyp  Einsatztyp  `json:"einsatzsubtyp"`
	Adresse     Adresse     `json:"adresse"`
	Feuerwehren Feuerwehren `json:"feuerwehren"`
}

type Feuerwehren []Feuerwehr

func (f *Feuerwehren) UnmarshalJSON(data []byte) error {
	// parse json '{"1337": {"feuerwehr": "Feuerwehr 1337"}, "420": {"feuerwehr": "Feuerwehr 420"}}' to Feuerwehr objects
	var feuerwehrenMap map[string]any
	err := json.Unmarshal(data, &feuerwehrenMap)
	if err != nil {
		return err
	}

	// Loop through each Feuerwehr object in the map and parse it
	for key, feuerwehrObject := range feuerwehrenMap {
		fwNr := key
		fwName := feuerwehrObject.(map[string]any)["feuerwehr"].(string)

		*f = append(*f, Feuerwehr{FwNr: fwNr, FwName: fwName})
	}

	return nil
}

type Koordinaten struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Einsatztyp struct {
	ID   string `json:"id"`
	Name string `json:"text"`
}

type Adresse struct {
	Bezirk     string `json:"text"`
	Ort        string `json:"emun"`
	Strasse    string `json:"efeanme"`
	Hausnummer string `json:"estnum"`
}

type Feuerwehr struct {
	FwNr   string `json:"fwnr"`
	FwName string `json:"fwname"`
}

type MyTime struct {
	time.Time
}

func (t *MyTime) UnmarshalJSON(data []byte) error {
	dateStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	date, err := time.Parse(time.RFC1123Z, dateStr)
	if err != nil {
		return err
	}

	t.Time = date
	return nil
}

func parseEinsatz(jsonString []byte) (Einsatz, error) {
	var einsatz Einsatz
	err := json.Unmarshal(jsonString, &einsatz)
	if err != nil {
		return einsatz, err
	}
	// sort Feuerwehren by FwNr
	sort.Slice(einsatz.Feuerwehren, func(i, j int) bool {
		return einsatz.Feuerwehren[i].FwNr < einsatz.Feuerwehren[j].FwNr
	})
	return einsatz, nil
}

func parseEinsaetze() []Einsatz {
	response, err := http.Get(jsonUrl)
	if err != nil {
		fmt.Println("Failed to get data from URL:", err)
		return nil
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil
	}

	// Parse the JSON data into a map
	var einsaetzeMap map[string]any
	err = json.Unmarshal(data, &einsaetzeMap)
	if err != nil {
		fmt.Println("Failed to parse JSON data:", err)
		return nil
	}

	if einsaetzeMap["einsaetze"] == nil {
		fmt.Println("No Einsätze found")
		return nil
	}

	var einsaetze []Einsatz

	einsaetzeMap = einsaetzeMap["einsaetze"].(map[string]any)

	// Loop through each Einsatz object in the map and parse it
	for _, einsatzObject := range einsaetzeMap {
		einsatzBytes, err := json.Marshal(einsatzObject.(map[string]any)["einsatz"].(map[string]any))
		if err != nil {
			fmt.Println("Failed to marshal einsatz object:", err)
			continue
		}

		einsatz, err := parseEinsatz(einsatzBytes)
		if err != nil {
			fmt.Println("Failed to parse einsatz object:", err)
			continue
		}

		einsaetze = append(einsaetze, einsatz)

		einsatzJson, err := json.Marshal(einsatz)
		if err != nil {
			fmt.Println("Failed to marshal einsatz object:", err)
			continue
		}
		fmt.Println(string(einsatzJson))
	}

	return einsaetze
}
