package mip

import (
	"fmt"
	"math"
)

// Constraint represents a linear constraint in the form of:
// a1*x1 + a2*x2 + ... + a_n*x_n {<=, >=, ==} b
type Constraint = constraint // export alias

// ConstraintType represents the sign between the linear expression and the right-hand side constant in a Constraint.
// Recall that a linear Constraint is in the form of a1*x1 + a2*x2 + ... + an*xn {<=, >=, ==} b
type ConstraintType string

const (
	LessThanOrEqual    = "<="
	Equal              = "=="
	GreaterThanOrEqual = ">="
)

// AddConstraintExpr adds a new Constraint to the Solver based on the given linear expression and Constraint type.
// For example if we want expression <= 5, we would call AddConstraintExpr(expr, LessThanOrEqual, 5.0)
func (s *Solver) AddConstraintExpr(e *LinearExpression, t ConstraintType, rhs float64) *Constraint {
	var c *Constraint

	switch t {
	case LessThanOrEqual:
		c = s.newConstraint(math.Inf(-1), rhs)
	case Equal, "=":
		c = s.newConstraint(rhs, rhs)
	case GreaterThanOrEqual:
		c = s.newConstraint(rhs, math.Inf(1))

	// In case "<" or ">" strings are directly passed to the function as ConstraintType
	case ">", "<":
		panic(fmt.Sprintf("Strict inequalities are not supported: %s", t))
	default:
		panic(fmt.Sprintf("Unknown cconstraint type: %s", t))
	}

	for v, weight := range e.terms {
		c.setCoefficient(v, weight)
	}

	return c
}
