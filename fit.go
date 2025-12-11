package main

import (
	"fmt"
	"math"
	"strings"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

const (
	lowerDegree int = 1
	upperDegree int = 9
)

// Representation of a two-dimentional data set with float numbers
type Data2D struct {
	x []float64
	y []float64
}

// Data2D constuctor. Returns error is x and y lenghts are mismatched.
func NewData2D(x, y []float64) (*Data2D, error) {
	if len(x) != len(y) {
		return nil, fmt.Errorf("x and y lenghts are mismatched, x lenght is %d, while y lenght is %d", len(x), len(y))
	}
	return &Data2D{x: x, y: y}, nil
}

func (d *Data2D) String() string {
	var output strings.Builder
	for i := 0; i < len(d.x); i++ {
		dataP := fmt.Sprintf("%d: x=%g, y=%g\n", i+1, d.x[i], d.y[i])
		output.WriteString(dataP)
	}
	res := output.String()
	// cut final LF
	return res[:len(res)-1]
}

// Representation of polynomial fitting of 2D dataset
type FitCurve struct {
	Polynom  Polynomial
	RSquared float64
}

func (f FitCurve) String() string {
	return fmt.Sprintf("f(x) = %v with R^2=%.5f", f.Polynom, f.RSquared)
}

// Fit a 2D dataset with polynomial curve of some degree between 1-9
func PolyFit(data Data2D, degree int) (*FitCurve, error) {
	if degree < lowerDegree || degree > upperDegree {
		return nil, fmt.Errorf("polynomial order should lay inside [1:9] interval")
	}
	if len(data.x) < degree+1 {
		return nil, fmt.Errorf("need at least %d points for degree %d polynomial", degree+1, degree)
	}

	fitted, err := qrFit(data, degree)
	if err != nil {
		return nil, err
	}

	predictedY := make([]float64, len(data.y))
	for i := range predictedY {
		predictedY[i] = fitted.Evaluate(data.x[i])
	}

	actualVSPredicted, err := NewData2D(data.y, predictedY)
	if err != nil {
		return nil, err
	}

	rSquared := rSquare(*actualVSPredicted)

	return &FitCurve{Polynom: *fitted, RSquared: rSquared}, nil
}

// Solving a system of linear equations via QR decomposition
// of Vandermonde matrix.
func qrFit(data Data2D, degree int) (*Polynomial, error) {
	n := len(data.x)

	vandermonde := mat.NewDense(n, degree+1, nil)
	for i := range n {
		for j := 0; j <= degree; j++ {
			vandermonde.Set(i, j, math.Pow(data.x[i], float64(j)))
		}
	}
	var qr mat.QR
	qr.Factorize(vandermonde)

	yVec := mat.NewVecDense(n, data.y)
	coeffs := mat.NewVecDense(degree+1, nil)

	err := qr.SolveVecTo(coeffs, false, yVec)
	if err != nil {
		return nil, fmt.Errorf("failed to solve: %v", err)
	}

	result := make([]float64, degree+1)
	for i := 0; i <= degree; i++ {
		result[i] = coeffs.AtVec(i)
	}

	return &Polynomial{coeffs: result}, nil
}

// RSquared calculates the coefficient of determination (R²)
// R² = 1 - (SS_res / SS_tot)
// where SS_res = Σ(y_i - y_pred_i)² and SS_tot = Σ(y_i - y_mean)²
func rSquare(actualVSPredicted Data2D) float64 {
	yActual := actualVSPredicted.x
	yPredicted := actualVSPredicted.y

	yMean := stat.Mean(yActual, nil)

	ssRes := 0.0
	for i := range yActual {
		residual := yActual[i] - yPredicted[i]
		ssRes += residual * residual
	}

	ssTot := 0.0
	for _, y := range yActual {
		diff := y - yMean
		ssTot += diff * diff
	}

	return 1.0 - (ssRes / ssTot)
}
