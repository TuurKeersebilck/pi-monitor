// Package roundutil provides the float-rounding helper shared by every
// backend package that formats a stat for display.
package roundutil

import "math"

// Round rounds val to the given number of decimal places.
func Round(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}
