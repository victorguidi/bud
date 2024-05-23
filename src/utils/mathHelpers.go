package utils

import (
	"log"
	"math"
)

func CosSimilarity(v1, v2 []float32) float64 {
	log.Println(len(v1), len(v2))
	if len(v1) != len(v2) {
		panic("Vectors should have the same size")
	}

	var a_b, v1_mag, v2_mag float64

	for i := range v1 {
		a_b += float64(v1[i]) * float64(v2[i])
		v1_mag += float64(v1[i]) * float64(v1[i])
		v2_mag += float64(v2[i]) * float64(v2[i])
	}

	v1_mag = math.Sqrt(v1_mag)
	v2_mag = math.Sqrt(v2_mag)

	return a_b / (v1_mag * v2_mag)
}

func Sigmoid(value float32) float64 {
	return 1 / (1 + math.Exp(-float64(value)))
}
