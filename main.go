package main

import (
	"Sun/sun"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type GlobalExport struct {
	//Time []int `json:"Time"`
	MinTemp   []float64 `json:"MinTemp"`
	Cloudness []int     `json:"Cloudness"`
	Qout      []float64 `json:"Losses"`
	Qd        []float64 `json:"Dynamic"`
	Qbs       []float64 `json:"Static"`
	Sunrise   []string  `json:"Sunrise"`
	Sunset    []string  `json:"Sunset"`
	//Duration

	BestAngle   float64 `json:"BestAngle"` // 90 - bestBeta
	BestAzimuth float64 `json:"BestAzimuth"`
}

const (
	colAr = 2.5 // collector active area
)

func main() {
	// if err := sun.GetYandex(); err != nil {
	// 	fmt.Print(err)
	// }
	// if err := sun.ExportInitialCSV(); err != nil {
	// 	fmt.Print(err)
	// }
	// if err := sun.GetSunTimes(); err != nil {
	// 	fmt.Print(err)
	// }

	// Sunrise, Sunset, err := sun.ExportSunTimes()
	// if err != nil {
	// 	fmt.Print(err)
	// }

	MinTemp, Cloudness, Qd, Qf, Qbsa, Qbsw, Qbss, err := sun.SolarIncome() // BestAngle, BestAzimuth,
	if err != nil {
		fmt.Print(err)
	}

	Qout, err := sun.GetLosses(MinTemp)
	if err != nil {
		fmt.Print(err)
	}

	fw, err := os.Create("data/CSVs/globalExport.csv")
	if err != nil {
		fmt.Print(err)
	}
	defer fw.Close()

	records := [][]string{
		{
			"Time",
			"MinTemp",
			"Cloudness",
			"Qout",
			"24h out",
			"Qd",
			"%",
			"Qf",
			"Qbs a",
			"Qbs w",
			"Qbs s",
			"% bs a/d",
			"% bs w/d",
			"% bs s/d",
			"% f/d",
			"Qhwp",
			"% hw",
			//"Sunrise",
			//"Sunset",
			// TODO: add duration and seconds.
		},
	}

	for i, temp := range MinTemp {
		records = append(records, []string{
			strconv.Itoa(i + 1),
			strconv.FormatFloat(temp, 'E', -2, 64),
			strconv.Itoa(Cloudness[i]),
			strconv.FormatFloat(Qout[i], 'E', -2, 64),
			strconv.FormatFloat(Qout[i]*24, 'E', -2, 64),
			strconv.FormatFloat(Qd[i]*colAr, 'E', -2, 64),
			strconv.FormatFloat(Qd[i]*colAr/(Qout[i]*24), 'E', -2, 64),
			strconv.FormatFloat(Qf[i]*colAr, 'E', -2, 64),
			strconv.FormatFloat(Qbsa[i]*colAr, 'E', -2, 64),
			strconv.FormatFloat(Qbsw[i]*colAr, 'E', -2, 64),
			strconv.FormatFloat(Qbss[i]*colAr, 'E', -2, 64),
			strconv.FormatFloat((Qd[i]-Qbsa[i])/Qd[i], 'E', -2, 64),
			strconv.FormatFloat((Qd[i]-Qbsw[i])/Qd[i], 'E', -2, 64),
			strconv.FormatFloat((Qd[i]-Qbss[i])/Qd[i], 'E', -2, 64),
			strconv.FormatFloat((Qd[i]-Qf[i])/Qd[i], 'E', -2, 64),
			strconv.FormatFloat(float64(4960), 'E', -2, 64),
			strconv.FormatFloat(Qd[i]*colAr/float64(4960), 'E', -2, 64),
			//Sunrise[i],
			//Sunset[i],
		})
	}

	w := csv.NewWriter(fw)

	if err := w.WriteAll(records); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

}
