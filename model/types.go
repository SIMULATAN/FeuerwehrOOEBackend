package model

import (
	"encoding/json"
	"time"
)

type Einsatz struct {
	ID          string       `json:"num1"`
	Einsatzort  string       `json:"einsatzort"`
	Koordinaten Koordinaten  `json:"wgs84"`
	Alarmstufe  json.Number  `json:"alarmstufe"`
	Startzeit   MyTime       `json:"startzeit"`
	Einsatztyp  Einsatztyp   `json:"einsatztyp"`
	Einsatzsubtyp Einsatztyp `json:"einsatzsubtyp"`
	Adresse     Adresse      `json:"adresse"`
	Feuerwehren Feuerwehren  `json:"feuerwehren"`
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
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
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

func (a Adresse) PrettyName() string {
	return a.Strasse + " " + a.Hausnummer + ", " + a.Ort + " (" + a.Bezirk + ")"
}

type Feuerwehr struct {
	FwNr   string `json:"fwnr"`
	FwName string `json:"fwname"`
}

type MyTime struct {
	time.Time
}
