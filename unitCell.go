package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Symmetry int

const (
	Triclinic Symmetry = iota
	Monoclinic
	Rhombohedral
	Hexagonal
	Orthorhombic
	Tetragonal
	Cubic
)

func (s Symmetry) String() string {
	return [...]string{"Triclinic", "Monoclinic", "Rhombohedral", "Hexagonal", "Orthorhombic", "Tetragonal", "Cubic"}[s]
}

type UnitCellData struct {
	CellSymmetry Symmetry
	T            []float64
	A            []float64
	B            []float64
	C            []float64
	Alpha        []float64
	Beta         []float64
	Gamma        []float64
}

func (u UnitCellData) String() string {
	var output strings.Builder
	output.WriteString(fmt.Sprintf("Symmetry: %s\n", u.CellSymmetry.String()))
	output.WriteString(fmt.Sprintf("T: %v\n", u.T))
	output.WriteString(fmt.Sprintf("a: %v\n", u.A))
	output.WriteString(fmt.Sprintf("b: %v\n", u.B))
	output.WriteString(fmt.Sprintf("c: %v\n", u.C))
	output.WriteString(fmt.Sprintf("alpha: %v\n", u.Alpha))
	output.WriteString(fmt.Sprintf("beta: %v\n", u.Beta))
	output.WriteString(fmt.Sprintf("gamma: %v", u.Gamma))
	return output.String()
}

func NewUnitCellData(sym Symmetry, t []float64, a []float64, params ...[]float64) (*UnitCellData, error) {
	n := len(t)
	degrees90 := fillSlice(n, 90.0)
	degrees120 := fillSlice(n, 120.0)

	ucd := &UnitCellData{
		CellSymmetry: sym,
		T:            t,
		A:            a,
	}

	switch sym {

	case Triclinic:
		if len(params) != 5 {
			return nil, fmt.Errorf("triclinic requires 5 additional params: b, c, alpha, beta, gamma")
		}
		ucd.B = params[0]
		ucd.C = params[1]
		ucd.Alpha = params[2]
		ucd.Beta = params[3]
		ucd.Gamma = params[4]

	case Monoclinic:
		if len(params) != 3 {
			return nil, fmt.Errorf("monoclinic requires 3 additional params: b, c, beta")
		}
		ucd.B = params[0]
		ucd.C = params[1]
		ucd.Alpha = slices.Clone(degrees90)
		ucd.Beta = params[2]
		ucd.Gamma = slices.Clone(degrees90)

	case Rhombohedral:
		if len(params) != 1 {
			return nil, fmt.Errorf("rhombohedral requires 1 additional param: alpha")
		}
		ucd.B = slices.Clone(ucd.A)
		ucd.C = slices.Clone(ucd.A)
		ucd.Alpha = params[0]
		ucd.Beta = slices.Clone(ucd.Alpha)
		ucd.Gamma = slices.Clone(ucd.Alpha)

	case Hexagonal:
		if len(params) != 1 {
			return nil, fmt.Errorf("hexagonal requires 1 additional param: c")
		}
		ucd.B = slices.Clone(ucd.A)
		ucd.C = params[0]
		ucd.Alpha = slices.Clone(degrees90)
		ucd.Beta = slices.Clone(degrees90)
		ucd.Gamma = slices.Clone(degrees120)

	case Orthorhombic:
		if len(params) != 2 {
			return nil, fmt.Errorf("orthorhombic requires 2 additional param: b, c")
		}
		ucd.B = params[0]
		ucd.C = params[1]
		ucd.Alpha = slices.Clone(degrees90)
		ucd.Beta = slices.Clone(degrees90)
		ucd.Gamma = slices.Clone(degrees90)

	case Tetragonal:
		if len(params) != 1 {
			return nil, fmt.Errorf("tetragonal requires 1 additional param: c")
		}
		ucd.B = slices.Clone(ucd.A)
		ucd.C = params[0]
		ucd.Alpha = slices.Clone(degrees90)
		ucd.Beta = slices.Clone(degrees90)
		ucd.Gamma = slices.Clone(degrees90)

	case Cubic:
		if len(params) != 0 {
			return nil, fmt.Errorf("cubic requires no additional params")
		}
		ucd.B = slices.Clone(ucd.A)
		ucd.C = slices.Clone(ucd.A)
		ucd.Alpha = slices.Clone(degrees90)
		ucd.Beta = slices.Clone(degrees90)
		ucd.Gamma = slices.Clone(degrees90)
	}

	err := validateDimensions(ucd.T, ucd.A, ucd.B, ucd.C, ucd.Alpha, ucd.Beta, ucd.Gamma)
	if err != nil {
		return nil, err
	}

	// turn degrees to radiants
	ucd.Alpha = degreesIntoRadians(ucd.Alpha)
	ucd.Beta = degreesIntoRadians(ucd.Beta)
	ucd.Gamma = degreesIntoRadians(ucd.Gamma)
	return ucd, nil
}

// Confirm that all of slices are the same length
func validateDimensions(params ...[]float64) error {
	n := len(params[0])
	for _, param := range params {
		if len(param) != n {
			return fmt.Errorf("some parameter with length %d, while it should be %d", len(param), n)
		}
	}
	return nil
}

// Fill a slice of known length with the same float
func fillSlice(length int, value float64) []float64 {
	res := make([]float64, length)
	for i := range length {
		res[i] = value
	}
	return res
}

// Turn degrees to radians by multiplying them by the factor of (pi/180)
func degreesIntoRadians(data []float64) []float64 {
	res := make([]float64, len(data))
	for i, datapoint := range data {
		res[i] = (math.Pi / 180) * datapoint
	}
	return res
}

// Read CSV file with data, parse it and turn into UnitCellData object
func CSVToUnitCell(filepath string, separator string, sym Symmetry) (*UnitCellData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	csvData, err := parseCSVDataIntoFloats(records, separator)
	if err != nil {
		return nil, err
	}

	rotated := pivotFloatSlices(csvData)
	ucd, err := NewUnitCellData(sym, rotated[0], rotated[1], rotated[2:]...)
	if err != nil {
		return nil, err
	}

	return ucd, nil
}

// Swap rows and cols of the table
func pivotFloatSlices(data [][]float64) [][]float64 {
	if len(data) == 0 {
		return [][]float64{}
	}

	rows := len(data)
	cols := len(data[0])

	pivoted := make([][]float64, cols)
	for i := range pivoted {
		pivoted[i] = make([]float64, rows)
	}

	for i := range rows {
		for j := range cols {
			pivoted[j][i] = data[i][j]
		}
	}

	return pivoted
}

// Parse string rows of floats with some separator
func parseCSVDataIntoFloats(records [][]string, sep string) ([][]float64, error) {
	csvData := [][]float64{}
	for _, record := range records {
		row := []float64{}
		for _, field := range record {
			parsed := strings.SplitSeq(field, sep)
			for parsedString := range parsed {
				if parsedString != "" {
					float, err := strconv.ParseFloat(parsedString, 64)
					if err != nil {
						return nil, err
					}
					row = append(row, float)
				}
			}
		}
		csvData = append(csvData, row)
	}
	return csvData, nil
}
