package sun

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	fi       = 54.99   // latitude; degrees
	Gsc      = 1368    // solar const W/m2
	radToDeg = 57.2958 // just ratio
)

var (
	// Tl is Turbidity coefficient; YEAR(location): 2.8
	Tl = []float64{2.0, 2.0, 2.0, 2.8, 3.8, 4.0, 4.2, 3.9, 3.2, 2.6, 2.0, 2.0}
)

type exportQin struct {
	Dynamic []float64 `json:"Dynamic"`
	Static  []float64 `json:"Static"`
}
type Position struct {
	azimuth, beta float64
}

type BestPosition struct {
	Qs            []float64
	Azimuth, Beta float64
}

func radians(degrees float64) float64 {
	return degrees / radToDeg
}

func degrees(radians float64) float64 {
	return radians * radToDeg
}

// finding best position by total income for anuual/winter/summer periods.
func findPosition(season string) ([]float64, float64, float64, error) {
	var bestQs []float64
	bestQsSum := float64(0)
	bestAzimuth := float64(0)
	bestBeta := float64(0)

	for azimuth := 50; azimuth <= 306; azimuth++ { // annual best – 115:33
		for beta := 0; beta <= 90; beta++ {
			Qs, _, _, err := getSolarIncome(Position{azimuth: float64(azimuth), beta: float64(beta)}, season)
			if err != nil {
				return nil, 0, 0, err
			}
			QsSum := float64(0)
			for _, Q := range Qs {
				QsSum += Q
			}
			bestQsSum = float64(0)
			for _, Q := range bestQs {
				bestQsSum += Q
			}
			if QsSum > bestQsSum {
				bestQs = Qs
				bestAzimuth = float64(azimuth)
				bestBeta = float64(beta)
				//fmt.Printf("azimuth: %v \n beta: %v \n sumQ: %v \n\n", bestAzimuth, bestBeta, QsSum)
			}
		}
	}

	bestQsSum = float64(0)
	for _, Q := range bestQs {
		bestQsSum += Q
	}
	text := "annual"
	if season != "" {
		text = season
	}
	fmt.Printf("FINAL %v: azimuth: %v \n beta: %v \n sumQ: %v \n", text, bestAzimuth, bestBeta, bestQsSum) // Q: %v ... , bestQs)

	return bestQs, bestAzimuth, bestBeta, nil
}

// SolarIncome calculating solar income depends on weather stats.
func SolarIncome() ([]float64, []int, []float64, []float64, []float64, []float64, []float64, error) {

	// dynamic system income calculation.
	Qd, MinTemp, Cloudness, err := getSolarIncome(Position{}, "")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	// fixed system income calculation (definitely choosen position).
	Qf, _, _, err := getSolarIncome(Position{azimuth: float64(245), beta: float64(5)}, "")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	bestQs, bestAzimuth, bestBeta, err := findPosition("")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	bestAnnual := BestPosition{
		Azimuth: bestAzimuth,
		Beta:    bestBeta,
		Qs:      bestQs,
	}
	bestQs, bestAzimuth, bestBeta, err = findPosition("winter")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	bestWinter := BestPosition{
		Azimuth: bestAzimuth,
		Beta:    bestBeta,
		Qs:      bestQs,
	}
	bestQs, bestAzimuth, bestBeta, err = findPosition("summer")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	bestSummer := BestPosition{
		Azimuth: bestAzimuth,
		Beta:    bestBeta,
		Qs:      bestQs,
	}

	// var w io.Writer
	// fw, err := os.OpenFile("data/exportQ.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	return nil, nil, nil, nil, 0, 0, err
	// }
	// defer fw.Close()
	// w = fw

	// var enc *json.Encoder
	// if w != nil {
	// 	enc = json.NewEncoder(w)
	// }

	// if enc != nil {
	// 	if err := enc.Encode(exportQin{Dynamic: Qd, Static: bestQs}); err != nil {
	// 		fmt.Print(err)
	// 	}
	// }

	return MinTemp, Cloudness, Qd, Qf, bestAnnual.Qs, bestWinter.Qs, bestSummer.Qs, nil
}

func getSolarIncome(position Position, season string) ([]float64, []float64, []int, error) {
	var Q []float64
	var positions [][]string

	fr, err := os.Open("data/weather.json")
	if err != nil {
		return nil, nil, nil, err
	}
	defer fr.Close()

	var weather []Weather
	if err := json.NewDecoder(fr).Decode(&weather); err != nil {
		return nil, nil, nil, err
	}

	f, err := os.Open("data/CSVs/SunEarthTools_1h.csv")
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		positions = append(positions, strings.Split(record[0], ";"))
	}

	var Cloudness []int
	var MinTemp []float64

	for i := 0; i < 365; i++ {

		switch season {
		case "winter":
			if i > 125 && i < 269 {
				Q = append(Q, float64(0))
				continue
			}
			break
		case "summer":
			if i < 125 || i > 269 {
				Q = append(Q, float64(0))
				continue
			}
			break
		default:
			break
		}

		dayN := weather[i].Time
		cloudness := weather[i].Cloudness
		Cloudness = append(Cloudness, int(weather[i].Cloudness))
		MinTemp = append(MinTemp, weather[i].MinDayTemp)

		// initials
		Qd := float64(0)
		azimuth := float64(0)
		beta := float64(0)

		for ii := 1; ii < 49; ii += 2 {
			if positions[i+1][ii+1] == "--" {
				continue
			}

			azimuth, err = strconv.ParseFloat(positions[i+1][ii+1], 64) // position; degrees
			if err != nil {
				return nil, nil, nil, err
			}
			beta, err = strconv.ParseFloat(positions[i+1][ii], 64) // elevation; degrees
			if err != nil {
				return nil, nil, nil, err
			}

			declination := -23.44 * math.Cos(radians(float64(360))/float64(365*(dayN+10)))
			// hour angle
			omega := math.Asin(math.Sin(radians(azimuth)) * math.Cos(radians(beta)) / math.Cos(radians(declination)))
			// zenith angle
			cosZenith := math.Cos(radians(float64(fi)))*math.Cos(radians(declination))*math.Cos(omega) + math.Sin(radians(float64(fi)))*math.Sin(radians(declination))
			if cosZenith < 0 || cosZenith > 90 {
				continue
			}

			intPart, _ := math.Modf(float64(i / 30))
			if intPart > 11 {
				intPart--
			}

			Gbn := Gsc * math.Pow(2.72, (-Tl[int(intPart)]/(0.9+9.4*cosZenith)))
			Gbt := float64(0)
			if (Position{}) == position {
				Gbt = Gbn
			} else {
				if math.Abs(position.azimuth-azimuth) < 90 {
					angle := degrees(math.Acos(cosZenith)) - (math.Abs(position.azimuth-azimuth) + math.Abs(position.beta+beta)) // ??
					Gbt = Gbn * math.Cos(radians(float64(angle)))
				} else {
					Gbt = 0
				}
			}
			if q := Gbt * float64(1-cloudness/100) / 0.8; q > float64(0) {
				Qd += q
				if q > 1367 {
					log.Fatalf("FATAL! q: %v \n Gbt: %v \n Gbn: %v \n cosZenith:: %v \n ", q, Gbt, Gbn, cosZenith)
				}
			}
		}

		Q = append(Q, Qd)
	}

	return Q, MinTemp, Cloudness, nil
}
