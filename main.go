package main

import "fmt"

func main() {
	degrees, err := NewDegress([]int{2, 2, 1, 1, 1, 1})
	if err != nil {
		fmt.Println(err)
	}
	ucd, err := CSVToUnitCell("Bi3NbTiO9.csv", ";", Orthorhombic)
	if err != nil {
		fmt.Println(err)
	}
	funcs, err := NewUnitCellFuncsOverT(*ucd, *degrees)
	fmt.Println(funcs)
	if err != nil {
		fmt.Println(err)
	}
	ders := NewDerivatives(*funcs)
	fmt.Println(ders)
}
