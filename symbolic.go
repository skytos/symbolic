package main

import (
	"fmt"
	"math"
)

type Expression interface {
	String() string
	Derivative(v Variable) Expression
	Evaluate(vals map[string]float64) float64
}

type Constant struct {
	value float64
}

type Variable struct {
	name string
}

type Sum struct {
	a, b Expression
}

type Product struct {
	a, b Expression
}

type Power struct {
	base, exponent Expression
}

type Sin struct {
	value Expression
}

type Cos struct {
	value Expression
}

var ONE = Constant{1}
var NEGATIVE_ONE = Constant{-1}
var ZERO = Constant{0}

func (self Constant) String() string {
	return fmt.Sprintf("%v", self.value)
}

func (self Constant) Derivative(v Variable) Expression {
	return ZERO
}

func (self Constant) Evaluate(_vals map[string]float64) float64 {
	return self.value
}

func (self Variable) String() string {
	return self.name
}

func (self Variable) Derivative(v Variable) Expression {
	if v.name == self.name {
		return ONE
	} else {
		return ZERO
	}
}

func (self Variable) Evaluate(vals map[string]float64) float64 {
	return vals[self.name]
}

func (self Sum) String() string {
	return fmt.Sprintf("(%v+%v)", self.a, self.b)
}

func (self Sum) Derivative(v Variable) Expression {
	return Sum{self.a.Derivative(v), self.b.Derivative(v)}
}

func (self Sum) Evaluate(vals map[string]float64) float64 {
	return self.a.Evaluate(vals) + self.b.Evaluate(vals)
}

func (self Product) String() string {
	return fmt.Sprintf("(%v*%v)", self.a, self.b)
}

func (self Product) Derivative(v Variable) Expression {
	return Sum{
		Product{self.a.Derivative(v), self.b},
		Product{self.a, self.b.Derivative(v)},
	}
}

func (self Product) Evaluate(vals map[string]float64) float64 {
	return self.a.Evaluate(vals) * self.b.Evaluate(vals)
}

func (self Power) String() string {
	return fmt.Sprintf("%v^%v", self.base, self.exponent)
}

func (self Power) Derivative(v Variable) Expression {
	return Product{
		self.exponent,
		Product{self.base.Derivative(v),
			Power{self.base,
				Sum{self.exponent, NEGATIVE_ONE}}}}
}

func (self Power) Evaluate(vals map[string]float64) float64 {
	return math.Pow(self.base.Evaluate(vals), self.exponent.Evaluate(vals))
}

func (self Sin) String() string {
	return fmt.Sprintf("sin(%v)", self.value)
}

func (self Sin) Derivative(v Variable) Expression {
	return Product{self.value.Derivative(v), Cos{self.value}}
}

func (self Sin) Evaluate(vals map[string]float64) float64 {
	return math.Sin(self.value.Evaluate(vals))
}
func (self Cos) String() string {
	return fmt.Sprintf("cos(%v)", self.value)
}

func (self Cos) Derivative(v Variable) Expression {
	return Product{negate(self.value.Derivative(v)), Sin{self.value}}
}

func (self Cos) Evaluate(vals map[string]float64) float64 {
	return math.Cos(self.value.Evaluate(vals))
}

func negate(e Expression) Expression {
	return Product{e, NEGATIVE_ONE}
}

func invert(e Expression) Expression {
	return Power{e, NEGATIVE_ONE}
}

func euler(e Expression, v Variable) Expression {
	return Sum{v, negate(Product{e, invert(e.Derivative(v))})}
}

func quadratic(a, b, c Expression) Expression {
	return Product{
		Sum{
			negate(b),
			negate(Power{
				Sum{
					Power{b, Constant{2}},
					Product{
						Constant{-4},
						Product{a, c},
					},
				},
				Constant{0.5},
			}),
		},
		invert(Product{Constant{2}, a}),
	}
}

func main() {
	h := Variable{"h"}
	s := Variable{"s"}
	a := Variable{"a"}

	time := quadratic(
		Constant{-9.8},
		Product{s, Sin{a}},
		h,
	)

	distance := Product{
		s,
		Product{
			Cos{a},
			time,
		},
	}

	derivativeOfDistance := distance.Derivative(a)

	for i, z := 0, 0.5; i < 10; i++ {
		fmt.Println(z)
		z = euler(derivativeOfDistance, a).Evaluate(
			map[string]float64{"h": 1.5, "s": 1.0, "a": z},
		)
	}
}
