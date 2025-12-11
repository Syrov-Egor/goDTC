package main

import (
	"fmt"
	"math"
	"strings"
)

type Polynomial struct {
	coeffs []float64
}

// Evaluate given polinomial function at given point x
// with Horner's method
func (p *Polynomial) Evaluate(x float64) float64 {
	if len(p.coeffs) == 0 {
		return 0
	}
	result := p.coeffs[len(p.coeffs)-1]
	for i := len(p.coeffs) - 2; i >= 0; i-- {
		result = result*x + p.coeffs[i]
	}

	return result
}

func (p Polynomial) String() string {
	if len(p.coeffs) == 0 {
		return "0"
	}

	var terms []string
	for i := len(p.coeffs) - 1; i >= 0; i-- {
		coeff := p.coeffs[i]
		if coeff == 0 {
			continue
		}

		var term string
		absCoeff := math.Abs(coeff)
		sign := ""
		if coeff > 0 && len(terms) > 0 {
			sign = " + "
		} else if coeff < 0 {
			sign = " - "
			if len(terms) == 0 {
				sign = " - "
			}
		}

		switch i {
		case 0:
			term = fmt.Sprintf("%s%.5f", sign, absCoeff)
		case 1:
			if absCoeff == 1 {
				term = fmt.Sprintf("%sx", sign)
			} else {
				term = fmt.Sprintf("%s%.gx", sign, absCoeff)
			}
		default:
			if absCoeff == 1 {
				term = fmt.Sprintf("%sx^%d", sign, i)
			} else {
				term = fmt.Sprintf("%s%.gx^%d", sign, absCoeff, i)
			}
		}

		terms = append(terms, term)
	}

	if len(terms) == 0 {
		return "0"
	}

	return strings.Join(terms, "")
}
