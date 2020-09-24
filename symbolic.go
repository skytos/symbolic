package main

import (
	"fmt"
	"math"
)

type expression interface {
	string() string
	derivative(v variable) expression
	evaluate(vals map[string]float64) float64
}

type constant struct {
	value float64
}

type variable struct {
	name string
}

type sum struct {
	a, b expression
}

type product struct {
	a, b expression
}

type power struct {
	base, exponent expression
}

type sin struct {
	value expression
}

type cos struct {
	value expression
}

var one = constant{1}
var negativeOne = constant{-1}
var zero = constant{0}

func (e constant) string() string {
	return fmt.Sprintf("%v", e.value)
}

func (e constant) derivative(v variable) expression {
	return zero
}

func (e constant) evaluate(_vals map[string]float64) float64 {
	return e.value
}

func (e variable) string() string {
	return e.name
}

func (e variable) derivative(v variable) expression {
	if v.name == e.name {
		return one
	}
	return zero
}

func (e variable) evaluate(vals map[string]float64) float64 {
	return vals[e.name]
}

func (e sum) string() string {
	return fmt.Sprintf("(%v+%v)", e.a, e.b)
}

func (e sum) derivative(v variable) expression {
	return sum{e.a.derivative(v), e.b.derivative(v)}
}

func (e sum) evaluate(vals map[string]float64) float64 {
	return e.a.evaluate(vals) + e.b.evaluate(vals)
}

func (e product) string() string {
	return fmt.Sprintf("(%v*%v)", e.a, e.b)
}

func (e product) derivative(v variable) expression {
	return sum{
		product{e.a.derivative(v), e.b},
		product{e.a, e.b.derivative(v)},
	}
}

func (e product) evaluate(vals map[string]float64) float64 {
	return e.a.evaluate(vals) * e.b.evaluate(vals)
}

func (e power) string() string {
	return fmt.Sprintf("%v^%v", e.base, e.exponent)
}

func (e power) derivative(v variable) expression {
	return product{
		e.exponent,
		product{e.base.derivative(v),
			power{e.base,
				sum{e.exponent, negativeOne}}}}
}

func (e power) evaluate(vals map[string]float64) float64 {
	return math.Pow(e.base.evaluate(vals), e.exponent.evaluate(vals))
}

func (e sin) string() string {
	return fmt.Sprintf("sin(%v)", e.value)
}

func (e sin) derivative(v variable) expression {
	return product{e.value.derivative(v), cos{e.value}}
}

func (e sin) evaluate(vals map[string]float64) float64 {
	return math.Sin(e.value.evaluate(vals))
}
func (e cos) string() string {
	return fmt.Sprintf("cos(%v)", e.value)
}

func (e cos) derivative(v variable) expression {
	return product{negate(e.value.derivative(v)), sin{e.value}}
}

func (e cos) evaluate(vals map[string]float64) float64 {
	return math.Cos(e.value.evaluate(vals))
}

func negate(e expression) expression {
	return product{e, negativeOne}
}

func invert(e expression) expression {
	return power{e, negativeOne}
}

func euler(e expression, v variable) expression {
	return sum{v, negate(product{e, invert(e.derivative(v))})}
}

func quadratic(a, b, c expression) expression {
	return product{
		sum{
			negate(b),
			negate(power{
				sum{
					power{b, constant{2}},
					product{
						constant{-4},
						product{a, c},
					},
				},
				constant{0.5},
			}),
		},
		invert(product{constant{2}, a}),
	}
}

func main() {
	h := variable{"h"}
	s := variable{"s"}
	a := variable{"a"}

	time := quadratic(
		constant{-9.8},
		product{s, sin{a}},
		h,
	)

	distance := product{
		s,
		product{
			cos{a},
			time,
		},
	}

	derivativeofdistance := distance.derivative(a)

	for i, z := 0, 0.5; i < 10; i++ {
		fmt.Println(z)
		z = euler(derivativeofdistance, a).evaluate(
			map[string]float64{"h": 1.5, "s": 1.0, "a": z},
		)
	}
}
