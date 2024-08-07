package statistics

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

//PrintStatistics
/*
Based on example from "Mastering Go 4th ed".
*/
func PrintStatistics(input string) {
	numbers, err := parseNumbers(input)
	if len(numbers) == 0 || err == nil {
		fmt.Println("No valid input provided, generating one")
		numbers = generateNumbers()
	}
	minF, maxF := calculateMinMax(numbers)
	mean := calculateMean(numbers)
	stdDev := calculateStdDev(numbers, mean)

	fmt.Printf("Number of values: %d\n", len(numbers))
	fmt.Printf("Min: %.5f\n", minF)
	fmt.Printf("Max: %.5f\n", maxF)
	fmt.Printf("Mean: %.5f\n", mean)
	fmt.Printf("StdDev: %.5f\n", stdDev)
}

func generateNumbers() []float64 {
	amount := rand.Intn(19) + 1
	numbers := make([]float64, amount)
	for i := 0; i < amount; i++ {
		numbers[i] = rand.Float64() * 10
	}
	return numbers
}

func parseNumbers(input string) ([]float64, error) {
	numbers := make([]float64, len(input))
	for i, numStr := range strings.Split(input, " ") {
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, err
		}
		numbers[i] = f
	}
	return numbers, nil
}

func calculateMinMax(numbers []float64) (float64, float64) {
	minF, maxF := math.MaxFloat64, math.SmallestNonzeroFloat64
	for _, n := range numbers {
		switch {
		case n < minF:
			minF = n
		case n > maxF:
			maxF = n
		}
	}
	return minF, maxF
}

func calculateMean(numbers []float64) float64 {
	var sum float64
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}

func calculateStdDev(numbers []float64, mean float64) float64 {
	var sum float64
	for _, number := range numbers {
		diff := number - mean
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(numbers)))
}
