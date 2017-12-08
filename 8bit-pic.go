package main

import "golang.org/x/tour/pic"
import "math"

func Pic(dx, dy int) [][]uint8 {
	ss := make([][]uint8, dy)
    for i := range ss {
        ss[i] = make([]uint8, dx)
    }
	
	for i := 0; i < dy; i++ {
        ss[i] = make([]uint8, dx)
        for j := 0; j < dx; j++ {
            ss[i][j] = uint8(math.Abs(float64(i^j)-float64(j*2)))
        }
	}
	
	return a
}

func main() {
	pic.Show(Pic)
}
