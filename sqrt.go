package main

import (
	"fmt"
	"math"
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

func main() {
	q := 99.0
	fmt.Println(Sqrt(q))
	fmt.Println(math.Sqrt(q))