package utils

import "math"

func CosSimilarity(v1, v2 []float64) float64 {
	a_b := func() float64 {
		var value float64
		for i := range v1 {
			value += v1[i] * v2[i]
		}
		return value
	}()

	v1_mag := func() float64 {
		var value float64
		for _, v := range v1 {
			value += v * v
		}
		return math.Sqrt(value)
	}()
	v2_mag := func() float64 {
		var value float64
		for _, v := range v2 {
			value += v * v
		}
		return math.Sqrt(value)
	}()

	return a_b / (v1_mag * v2_mag)
}
