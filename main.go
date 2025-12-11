package main

import (
	"fmt"
)

func main() {
	x := []float64{298, 373, 473, 573, 673, 773, 873, 973, 1073, 1173, 1273}
	y := []float64{7.749, 7.75, 7.755, 7.749, 7.767, 7.767, 7.775, 7.775, 7.782, 7.784, 7.792}
	data, err := NewData2D(x, y)
	if err != nil {
		fmt.Println(err)
	}
	fitRes, err := PolyFit(*data, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fitRes)
}
