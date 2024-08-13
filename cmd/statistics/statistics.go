package statistics

import (
	"cmp"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

type statistics struct {
	NumberOfValues float64
	Min            float64
	Max            float64
	Mean           float64
	StdDev         float64
}

var sortedByMean = func(s1, s2 *statistics) int {
	return cmp.Compare(s1.Mean, s2.Mean)
}

func (s *statistics) String() string {
	return fmt.Sprintf("Num: %.2f | Min: %.2f | Max: %.2f | Mean: %.2f | StdDev: %.2f",
		s.NumberOfValues, s.Min, s.Max, s.Mean, s.StdDev)
}

//RunWriteExample
/*
Based on example from "Mastering Go 4th ed".
*/
func RunWriteExample(input string) {
	numbers, err := parseNumbers(input)
	if len(numbers) == 0 || err == nil {
		fmt.Println("No valid input provided, generating one")
		numbers = generateNumbers()
	}

	statistics := calculateStatistics(numbers)
	err = storeAsCsv(statistics)
	if err != nil {
		fmt.Println("Unable to write to file ", err)
	}
}

func RunReadExample() {
	statistics, err := readCsv()
	if err != nil {
		fmt.Println("Unable to read the file ", err)
		os.Exit(1)
	}
	slices.SortFunc(statistics, sortedByMean)
	for _, statistic := range statistics {
		fmt.Printf("%v\n", statistic)
	}
}

func generateNumbers() []float64 {
	amount := rand.Intn(19) + 1
	numbers := make([]float64, amount)
	for i := 0; i < amount; i++ {
		numbers[i] = rand.Float64() * 10
	}
	return numbers
}

func calculateStatistics(numbers []float64) *statistics {
	minF, maxF := calculateMinMax(numbers)
	mean := calculateMean(numbers)
	return &statistics{
		NumberOfValues: float64(len(numbers)),
		Min:            minF,
		Max:            maxF,
		Mean:           mean,
		StdDev:         calculateStdDev(numbers, mean),
	}
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
