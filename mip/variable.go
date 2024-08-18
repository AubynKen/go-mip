package mip

import "C"

// Variable represents a decision variable in the optimization problem.
// This is an export alias
type Variable = variable

// VarInt creates and returns a new integer Variable
func (s *Solver) VarInt(name string, lowerBound, upperBound int) *Variable {
	return s.newVariable(name, float64(lowerBound), float64(upperBound), 1)
}

// VarFloat creates and returns a new continuous Variable
func (s *Solver) VarFloat(name string, lowerBound, upperBound float64) *Variable {
	return s.newVariable(name, lowerBound, upperBound, 0)
}

// VarBool creates and returns a new decision/binary Variable
func (s *Solver) VarBool(name string) *Variable {
	return s.newVariable(name, 0., 1., 1)
}

// Name returns the name of the variable.
func (v *Variable) Name() string { return v.name() }

// Value returns the value of the variable in the solution after optimization.
func (v *Variable) Value() float64 { return v.solutionValue() }
