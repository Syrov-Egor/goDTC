package main

import "fmt"

func main() {
	ucd, err := CSVToUnitCell("Bi3NbTiO9.csv", Orthorhombic)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ucd)
}
