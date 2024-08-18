package mip

// LinearExpression represents a linear expression, in the form of:
// a1*x1 + a2*x2 + ... + a_n*x_n
// where a1, a2, ..., a_n are coefficients, and x1, x2, ..., x_n are variables.
type LinearExpression struct {
	terms map[*Variable]float64 // Variable to their corresponding coefficients
}

// NewLinearExpression creates an empty linear expression.
func NewLinearExpression() *LinearExpression {
	return &LinearExpression{
		terms: make(map[*Variable]float64),
	}
}

// AddTerm adds a new weighted term to the linear expression.
// i.e. (2 * x + 3 * y).AddTerm(z, 4) => 2 * x + 3 * y + 4 * z
func (e *LinearExpression) AddTerm(v *Variable, weight float64) {
	e.terms[v] = e.terms[v] + weight // if v is not in the map, e.terms[v] will be initially zero
}

func (e *LinearExpression) AddVar(v *Variable) {
	e.AddTerm(v, 1)
}

func (e *LinearExpression) AddExpr(other *LinearExpression) {
	for exprVariable, weight := range other.terms {
		e.AddTerm(exprVariable, weight)
	}
}
