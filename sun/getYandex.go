package sun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Weather struct {
	Time int64 `json:"time"`

	Cloudness float64 `json:"cl"`

	WindDir      string  `json:"wind_dir"`
	MaxWindSpeed float64 `json:"max_ws"`
	MinWindSpeed float64 `json:"min_ws"`

	MaxDayTemp float64 `json:"max_day_t"`
	MinDayTemp float64 `json:"min_night_t"`
	WaterTemp  float64 `json:"twater"`

	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`

	Icon       string  `json:"icon"`
	Condition  string  `json:"condition"`
	PrecChance float64 `json:"prec_chance"`
}

type Location struct {
	LocationId int    `json:"geoid"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
}

// GetYandex getting initial annual weather stats using yandex API.
func GetYandex() error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.weather.yandex.ru/v1/locations/66/longterm", nil) // 66 – Omsk; "https://api.weather.yandex.ru/v1/locations?lang=ru_RU"

	req.Header.Add("X-Yandex-API-Key", "70814571-b237-479a-a971-9e4aea834a35")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b := &bytes.Buffer{}
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return err
	}

	var data []Weather // var data []Location
	if err := json.NewDecoder(b).Decode(&data); err != nil {
		return err
	}

	var w io.Writer
	f, err := os.OpenFile("data/weather.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // locations
	if err != nil {
		return err
	}
	defer f.Close()
	w = f

	var enc *json.Encoder
	if w != nil {
		enc = json.NewEncoder(w)
	}

	if enc != nil {
		if err := enc.Encode(data); err != nil {
			fmt.Print("Unable to encode tickers to JSON log")
		}
	}

	return nil
}
