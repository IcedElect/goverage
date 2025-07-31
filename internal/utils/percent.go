package utils

import "math"

func Percent(covered, total int64) float64 {
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	
	percent := 100.0 * float64(covered) / float64(total)

	return math.Round(percent*100) / 100 // Round to two decimal places.
}