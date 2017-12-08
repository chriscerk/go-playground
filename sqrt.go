package main

import (
	"fmt"
	"math"
	"time"
	"math/rand"
)

func Sqrt(x float64) float64 {		
	return SqrtHelper(x, 1.0, 1.0)
}

func SqrtHelper(x, z, i float64) float64 {
	if i > x*10 {
		return z
	}
		
	return SqrtHelper(x, z-((z*z - x) / (2*z)), i+1)
}

func CompareResults(q, mySqrt, mathSqrt float64) {
	defer fmt.Println("\n")
	defer fmt.Printf("\n math.Sqrt: %d", mathSqrt)
	defer fmt.Printf("\n My Sqrt:   %d", mySqrt)

	fmt.Printf("\n For the Square Root of %s", int(q))
	fmt.Printf("\n Difference Between: %d", math.Abs(mySqrt - mathSqrt))
} 

func main() {
	rand.Seed(time.Now().UnixNano())

	q := float64(rand.Intn(10000))
	mySqrt := Sqrt(q)
	mathSqrt := math.Sqrt(q)

	CompareResults(q, mySqrt, mathSqrt)
}