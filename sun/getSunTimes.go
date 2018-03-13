package sun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type SunTime struct {
	Rise string `json:"sunrise"`
	Set  string `json:"sunset"`
	//"solar_noon"
	Length string `json:"day_length"`
	// "civil_twilight_begin"
	// "civil_twilight_end"
	// "nautical_twilight_begin"
	// "nautical_twilight_end"
	// "astronomical_twilight_begin"
	// "astronomical_twilight_end"
}

type SunTimeWithDate struct {
	Rise   string `json:"sunrise"`
	Set    string `json:"sunset"`
	Length string `json:"day_length"`
	Date   string `json:"date"`
}

type SunTimeObj struct {
	Results SunTime `json:"results"`
	Status  string  `json:"status"`
}

// GetSunTimes getting day light info by location using external API.
func GetSunTimes() error {

	var dataArray []SunTimeWithDate

	for i := 1; i < 13; i++ {

		for ii := 1; ii < 32; ii++ {

			link := "https://api.sunrise-sunset.org/json?lat=54.9900000&lng=73.3700000&date=2017-" + strconv.Itoa(i) + "-" + strconv.Itoa(ii)

			client := &http.Client{}
			req, err := http.NewRequest("GET", link, nil)

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
			if b.String() != "Array" {

				var data SunTimeObj
				if err := json.NewDecoder(b).Decode(&data); err != nil {
					return err
				}

				dataArray = append(dataArray, SunTimeWithDate{Rise: data.Results.Rise, Set: data.Results.Set, Length: data.Results.Length, Date: "2017-" + strconv.Itoa(i) + "-" + strconv.Itoa(ii)})

			}
		}
	}

	var w io.Writer
	f, err := os.OpenFile("data/sunTimes.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
		if err := enc.Encode(dataArray); err != nil {
			fmt.Print("Unable to encode tickers to JSON log")
		}
	}

	return nil
}
