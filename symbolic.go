package main

import (
	"fmt"
	"math"
)

type expression interface {
	String() string
	derivative(v variable) expression
	evaluate(vals map[string]float64) float64
	simplify() expression
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

var one = constant{1}
var negativeOne = constant{-1}
var zero = constant{0}

func (e constant) String() string {
	return fmt.Sprintf("%v", e.value)
}

func (e constant) derivative(v variable) expression {
	return zero
}

func (e constant) evaluate(_vals map[string]float64) float64 {
	return e.value
}

func (e constant) simplify() expression {
	return e
}

func (e variable) String() string {
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

func (e variable) simplify() expression{
	return e
}

func (e sum) String() string {
	return fmt.Sprintf("(%v + %v)", e.a, e.b)
}

func (e sum) derivative(v variable) expression {
	return sum{e.a.derivative(v), e.b.derivative(v)}
}

func (e sum) evaluate(vals map[string]float64) float64 {
	return e.a.evaluate(vals) + e.b.evaluate(vals)
}

var ZERO = constant{0.0}
var ONE = constant{1.0}

func (e sum) simplify() expression {
	a := e.a.simplify()
	b := e.b.simplify()

	if a == ZERO {
		return b
	} else if b == ZERO {
		return a
	} else {
		ac, a_ok := a.(constant)
		bc, b_ok := b.(constant)
		if a_ok && b_ok {
			return constant{ac.value + bc.value}
		} else {
			return sum{a, b}
		}
	}
}

func (e product) String() string {
	return fmt.Sprintf("(%v * %v)", e.a, e.b)
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

func (e product) simplify() expression {
	a := e.a.simplify()
	b := e.b.simplify()

	if a == ZERO {
		return a
	} else if b == ZERO {
		return b
	} else if a == ONE {
		return b
	} else if b == ONE {
		return a
	} else {
		ac, a_ok := a.(constant)
		bc, b_ok := b.(constant)
		if a_ok && b_ok {
			return constant{ac.value * bc.value}
		} else {
			return product{a, b}
		}
	}
}

func (e power) String() string {
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

func (e power) simplify() expression {
	base := e.base.simplify()
	exponent := e.exponent.simplify()

	if exponent == ZERO {
		return ONE
	} else if exponent == ONE {
		return base
	} else if base == ZERO || base == ONE{
		return base
	} else {
		basec, base_ok := base.(constant)
		exponentc, exponent_ok := exponent.(constant)
		if base_ok && exponent_ok {
			return constant{math.Pow(basec.value, exponentc.value)}
		} else {
			return power{base, exponent}
		}
	}
}

func negate(e expression) expression {
	return product{e, negativeOne}
}

func invert(e expression) expression {
	return power{e, negativeOne}
}

func add(e1, e2 expression) expression {
	return sum{e1,e2}
}

func sub(e1, e2 expression) expression {
	return sum{e1,negate(e2)}
}

func mul(e1, e2 expression) expression {
	return product{e1,e2}
}

func div(e1, e2 expression) expression {
	return product{e1,invert(e2)}
}

func pow(e1, e2 expression) expression {
	return power{e1,e2}
}

func main() {
	two := constant{2.0}
	three := constant{3.0}
	four := constant{4.0}

	m1 := variable{"m1"}
	m2 := variable{"m2"}
	m3 := variable{"m3"}
	m4 := variable{"m4"}
	b3 := variable{"b3"}
	
	b4 := div(mul(b3,m4), m3)
	l := negate(div(b3,m3))
	h1 := div(mul(m1,b3), sub(m1, m3))
	h2 := div(mul(m2,b3), sub(m2, m3))
	h3 := div(mul(m2,b4), sub(m2, m4))
	h4 := div(mul(m1,b4), sub(m1, m4))
	a := div(mul(l, h2), two)
	b := sub(div(mul(l, h1), two), a)
	c := sub(div(mul(l, h3), two), a)
	d := sub(sub(sub(div(mul(l, h4), two), a), b), c)

	ea := pow(sub(four, a), two)
	eb := pow(sub(three, b), two)
	ec := pow(sub(two, c), two)

	e := add(ea, add(eb, ec)).simplify()

	gm1 := e.derivative(m1).simplify()
	gm2 := e.derivative(m2).simplify()
	gm3 := e.derivative(m3).simplify()
	gm4 := e.derivative(m4).simplify()
	gb3 := e.derivative(b3).simplify()

	vars := map[string]float64{
		"m1": 1.0,
		"m2": 0.5,
		"m3": -0.5,
		"m4": -1.0,
		"b3": 2.0,
	}
	for i := 0; i < 15000; i++ {
		// fmt.Println(math.Log(e.evaluate(vars)))
		delta := 0.0007
		dm1 := -delta * gm1.evaluate(vars)
		dm2 := -delta * gm2.evaluate(vars)
		dm3 := -delta * gm3.evaluate(vars)
		dm4 := -delta * gm4.evaluate(vars)
		db3 := -delta * gb3.evaluate(vars)

		vars["m1"] += dm1
		vars["m2"] += dm2
		vars["m3"] += dm3
		vars["m4"] += dm4
		vars["b3"] += db3

	}
	fmt.Println(a.evaluate(vars), b.evaluate(vars), c.evaluate(vars), d.evaluate(vars))
}
