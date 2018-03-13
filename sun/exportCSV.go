package sun

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

// ExportInitialCSV exporting initial data from weather provider (yandex).
func ExportInitialCSV() error {
	fr, err := os.Open("data/weather_110117.json")
	if err != nil {
		return err
	}
	defer fr.Close()

	fw, err := os.Create("data/CSVs/weather.csv")
	if err != nil {
		return err
	}
	defer fw.Close()

	var data []Weather
	if err := json.NewDecoder(fr).Decode(&data); err != nil {
		return err
	}

	records := [][]string{
		{
			"Time",
			"Cloudness",
			//"WindDir",
			"MaxWindSpeed",
			//"MinWindSpeed",
			//"MaxDayTemp",
			"MinDayTemp",
			"WaterTemp",
			//"Condition",
			//"PrecChance",
		},
	}

	for _, day := range data {
		records = append(records, []string{
			strconv.FormatInt(day.Time, 10),
			strconv.FormatFloat(day.Cloudness, 'f', 3, 64),
			//day.WindDir,
			strconv.FormatFloat(day.MaxWindSpeed, 'f', 3, 64),
			//strconv.FormatFloat(day.MinWindSpeed, 'f', 3, 64),
			//strconv.FormatFloat(day.MaxDayTemp, 'f', 3, 64),
			strconv.FormatFloat(day.MinDayTemp, 'f', 3, 64),
			strconv.FormatFloat(day.WaterTemp, 'f', 3, 64),
			//day.Condition,
			//strconv.FormatFloat(day.PrecChance, 'f', 3, 64),
		})
	}

	w := csv.NewWriter(fw)

	if err := w.WriteAll(records); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	return nil
}
