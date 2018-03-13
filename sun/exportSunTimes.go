package sun

import (
	"encoding/json"
	"os"
)

// ExportSunTimes exporting sunrise, sunset times and duration of the light day.
func ExportSunTimes() ([]string, []string, error) {

	fr, err := os.Open("data/sunTimes.json")
	if err != nil {
		return nil, nil, err
	}
	defer fr.Close()

	fw, err := os.Create("data/CSVs/sunTimes.csv")
	if err != nil {
		return nil, nil, err
	}
	defer fw.Close()

	var data []SunTimeWithDate
	if err := json.NewDecoder(fr).Decode(&data); err != nil {
		return nil, nil, err
	}

	records := [][]string{
		{
			"Rise",
			"Set",
			"Length",
			"Date",
		},
	}

	for _, day := range data {
		records = append(records, []string{
			day.Rise,
			day.Set,
			day.Length,
			day.Date,
		})
	}

	var Sunrises, Sunsets []string

	for i, rec := range records {
		if i != 0 {
			Sunrises = append(Sunrises, rec[0])
			Sunsets = append(Sunsets, rec[1])
			// TODO: calculate duration using date parsing
		}
	}

	// w := csv.NewWriter(fw)

	// if err := w.WriteAll(records); err != nil {
	// 	log.Fatalln("error writing record to csv:", err)
	// }

	return Sunrises, Sunsets, nil
}
