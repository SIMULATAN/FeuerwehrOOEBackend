package model

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const LaufendUrl string = "https://cf-einsaetze.ooelfv.at/webext2/rss/json_taeglich.txt"

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

func (t *MyNullableTime) UnmarshalJSON(data []byte) error {
	dateStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	// ongoing Einsatz
	if (len(dateStr) == 0) {
		t.Time = nil
		return nil
	}

	date, err := time.Parse(time.RFC1123Z, dateStr)
	if err != nil {
		return err
	}

	t.Time = &date
	return nil
}

func (s MyNullableTime) MarshalJSON() ([]byte, error) {
	if s.Time != nil {
		return json.Marshal(s.Time)
	}
	return []byte(`null`), nil
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

func ParseLaufendeEinsaetze() []Einsatz {
	response, err := http.Get(LaufendUrl)
	if err != nil {
		log.Println("Failed to get data from URL:", err)
		return nil
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil
	}

	// Parse the JSON data into a map
	var einsaetzeMap map[string]any
	err = json.Unmarshal(data, &einsaetzeMap)
	if err != nil {
		log.Println("Failed to parse JSON data:", err)
		return nil
	}

	if einsaetzeMap["einsaetze"] == nil {
		log.Println("No einsaetze object found")
		return nil
	}

	var einsaetze []Einsatz

	einsaetzeMap = einsaetzeMap["einsaetze"].(map[string]any)

	// Loop through each Einsatz object in the map and parse it
	for _, einsatzObject := range einsaetzeMap {
		einsatzBytes, err := json.Marshal(einsatzObject.(map[string]any)["einsatz"].(map[string]any))
		if err != nil {
			log.Println("Failed to marshal einsatz object:", err)
			continue
		}

		einsatz, err := parseEinsatz(einsatzBytes)
		if err != nil {
			log.Println("Failed to parse einsatz object:", err)
			continue
		}

		einsaetze = append(einsaetze, einsatz)

		einsatzJson, err := json.Marshal(einsatz)
		if err != nil {
			log.Println("Failed to marshal einsatz object:", err)
			continue
		}

		log.Println("Parsed Einsatz:", string(einsatzJson))
	}

	return einsaetze
}
