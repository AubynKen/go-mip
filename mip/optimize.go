package mip

type OptimizationType int

const (
	Maximize OptimizationType = iota
	Minimize
)

// SetOptimizationObjective sets the objective function of the Solver to the given linear expression and optimization type.
func (s *Solver) SetOptimizationObjective(le *LinearExpression, tp OptimizationType) {
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
