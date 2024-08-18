package mip

import "C"

type OptimizationType int

const (
	Maximize OptimizationType = iota
	Minimize
)

// SetObjective sets the objective function of the Solver to the given linear expression and optimization type.
func (s *Solver) SetObjective(le *LinearExpression, tp OptimizationType) {
	svr := s

	for v, coefficient := range le.terms {
		svr.setObjectiveCoefficient(v, coefficient)
	}

	switch tp {
	case Maximize:
		svr.setMaximization()
	case Minimize:
		svr.setMinimization()
	}
}

// ObjectiveValue returns the current best objective value found by the solver.
func (s *Solver) ObjectiveValue() float64 {
	return s.objectiveValue()
}

// BestBound returns what's currently the best bound.
// for example, if the problem is being minimized, and the best bound is 100,
// then the theoretical optimal objective is at least 100.
// This can be used to evaluate the quality of the solution.
func (s *Solver) BestBound() float64 {
	return s.getBestBound()
}

// Gap returns the relative gap between the best integer solution found and the best bound.
func (s *Solver) Gap() float64 {
	return s.getGap()
}

// ResultStatus represents the status of the optimization result.
type ResultStatus int

const (
	Optimal ResultStatus = iota
	Feasible
	Infeasible
	Unbounded
	Abnormal
	ModelInvalid
	NotSolved
)
