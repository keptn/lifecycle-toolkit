package aggregation

import (
	"math"
	"sort"
)

func CalculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func CalculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}

func CalculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	// Sort the values
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	sort.Float64s(sortedValues)

	// Calculate the median
	middle := len(sortedValues) / 2
	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}
	return sortedValues[middle]
}

func CalculateAverage(values []float64) float64 {
	sum := 0.0

	for _, value := range values {
		sum += value
	}
	if len(values) > 0 {
		return sum / float64(len(values))
	}

	return 0.0
}

func CalculatePercentile(values sort.Float64Slice, perc float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	// Calculate the index for the requested percentile
	i := math.Ceil(float64(len(values)) * perc / 100)

	return values[int(i-1)]
}
