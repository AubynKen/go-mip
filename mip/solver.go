package mip

import "C"
import (
	"fmt"
	"time"
)

// Solver represents the optimization problem to be solved.
// This is an export alias
type Solver = solver

const (
	SCIP = "SCIP"
	CBC  = "CBC"
)

// NewSolver creates and returns a new Solver of the given type.
func NewSolver(solverType string) (*Solver, error) {

	if solverType != SCIP && solverType != CBC {
		return nil, fmt.Errorf("unsupported solver type")
	}

	solver := createSolver(solverType)
	if solver == nil {
		return nil, fmt.Errorf("failed to create Solver")
	}

	return solver, nil
}

// ReleaseResources frees up the memory in the C heap allocated for the Solver.
func (s *Solver) ReleaseResources() {
	s.delete()
}

// Solve attempts to solve the optimization problem within the given time limit.
// It returns a SolveResult containing the solution status, objective value, best bound, and gap.
func (s *Solver) Solve(timeLimit time.Duration) (isOptimal bool, err error) {
	if timeLimit > 0 {
		s.setTimeLimit(timeLimit.Milliseconds())
	}

	status := ResultStatus(s.solve())

	switch status {
	case Optimal:
		return true, nil
	case Feasible:
		return false, nil
	case Infeasible:
		return false, fmt.Errorf("the problem is infeasible")
	case Unbounded:
		return false, fmt.Errorf("the problem is unbounded")
	case Abnormal:
		return false, fmt.Errorf("the solver encountered an abnormal situation")
	case ModelInvalid:
		return false, fmt.Errorf("the model is invalid")
	case NotSolved:
		return false, fmt.Errorf("the problem was not solved")
	default:
		return false, fmt.Errorf("unknown result status")
	}
}
