package sun

// GetLosses calculating losses through the walls.
func GetLosses(MinTemps []float64) ([]float64, error) {
	var Qout []float64

	for _, temp := range MinTemps {
		k := 0.149844784
		Q := float64(50) * k * (float64(20) - temp)
		Qout = append(Qout, Q)
	}

	return Qout, nil
}
