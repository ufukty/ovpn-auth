package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

func Sum[N constraints.Integer | constraints.Float](s []N) N {
	total := N(0)
	for _, item := range s {
		total += item
	}
	return total
}

func Avg[N constraints.Integer | constraints.Float](s []N) float64 {
	return float64(Sum(s)) / float64(len(s))
}

func StandardDeviationForComparison[N constraints.Integer | constraints.Float](s []N) float64 {
	avg := Avg(s)
	total := 0.0
	for i := 0; i < len(s); i++ {
		total += math.Pow(float64(s[i])-avg, 2)
	}
	return total
}

func StandardDeviation[N constraints.Integer | constraints.Float](s []N) float64 {
	return math.Sqrt(StandardDeviationForComparison(s) / float64(len(s)))
}
