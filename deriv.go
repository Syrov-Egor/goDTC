package main

import (
	"fmt"
	"strings"
)

type Degrees struct {
	A     int
	B     int
	C     int
	Alpha int
	Beta  int
	Gamma int
}

func NewDegress(degrees []int) (*Degrees, error) {
	if len(degrees) != 6 {
		return nil, fmt.Errorf("there should be 6 degrees for polynomial functions, got %d", len(degrees))
	}
	for _, degree := range degrees {
		if degree < 1 {
			return nil, fmt.Errorf("there is a degree <1")
		}
	}
	return &Degrees{A: degrees[0], B: degrees[1], C: degrees[2], Alpha: degrees[3], Beta: degrees[4], Gamma: degrees[5]}, nil
}

type UnitCellFuncsOverT struct {
	A     FitCurve
	B     FitCurve
	C     FitCurve
	Alpha FitCurve
	Beta  FitCurve
	Gamma FitCurve
}

func NewUnitCellFuncsOverT(data UnitCellData, degrees Degrees) (*UnitCellFuncsOverT, error) {
	a, err := fitData(data.T, data.A, degrees.A)
	if err != nil {
		return nil, err
	}
	b, err := fitData(data.T, data.B, degrees.B)
	if err != nil {
		return nil, err
	}
	c, err := fitData(data.T, data.C, degrees.C)
	if err != nil {
		return nil, err
	}
	alpha, err := fitData(data.T, data.Alpha, degrees.Alpha)
	if err != nil {
		return nil, err
	}
	beta, err := fitData(data.T, data.Beta, degrees.Beta)
	if err != nil {
		return nil, err
	}
	gamma, err := fitData(data.T, data.Gamma, degrees.Gamma)
	if err != nil {
		return nil, err
	}
	return &UnitCellFuncsOverT{A: *a, B: *b, C: *c, Alpha: *alpha, Beta: *beta, Gamma: *gamma}, nil
}

func (u UnitCellFuncsOverT) String() string {
	var output strings.Builder
	output.WriteString(fmt.Sprintf("a: %v\n", u.A))
	output.WriteString(fmt.Sprintf("b: %v\n", u.B))
	output.WriteString(fmt.Sprintf("c: %v\n", u.C))
	output.WriteString(fmt.Sprintf("alpha: %v\n", u.Alpha))
	output.WriteString(fmt.Sprintf("beta: %v\n", u.Beta))
	output.WriteString(fmt.Sprintf("gamma: %v", u.Gamma))
	return output.String()
}

func fitData(Tdata []float64, Pdata []float64, degree int) (*FitCurve, error) {
	aData, err := NewData2D(Tdata, Pdata)
	if err != nil {
		return nil, err
	}
	a, err := PolyFit(*aData, degree)
	if err != nil {
		return nil, err
	}
	return a, nil
}

type Derivatives struct {
	AdT     Polynomial
	BdT     Polynomial
	CdT     Polynomial
	AlphadT Polynomial
	BetadT  Polynomial
	GammadT Polynomial
}

func NewDerivatives(funcs UnitCellFuncsOverT) *Derivatives {
	return &Derivatives{
		AdT:     *funcs.A.Polynom.Derivative(),
		BdT:     *funcs.B.Polynom.Derivative(),
		CdT:     *funcs.C.Polynom.Derivative(),
		AlphadT: *funcs.Alpha.Polynom.Derivative(),
		BetadT:  *funcs.Beta.Polynom.Derivative(),
		GammadT: *funcs.Gamma.Polynom.Derivative(),
	}
}

func (d Derivatives) String() string {
	var output strings.Builder
	output.WriteString(fmt.Sprintf("da/dT: %v\n", d.AdT))
	output.WriteString(fmt.Sprintf("db/dT: %v\n", d.BdT))
	output.WriteString(fmt.Sprintf("dc/dT: %v\n", d.CdT))
	output.WriteString(fmt.Sprintf("dalpha/dt: %v\n", d.AlphadT))
	output.WriteString(fmt.Sprintf("dbeta/dt: %v\n", d.BetadT))
	output.WriteString(fmt.Sprintf("dgamma/dt: %v", d.GammadT))
	return output.String()
}
