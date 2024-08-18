package examples

import (
	"fmt"
	"log"

	"mip-bridge-without-swig/mip"
)

func knapsackProblem() {
	// Create a new solver
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		log.Fatalf("Error creating solver: %v", err)
	}
	defer solver.Delete()

	// Define the knapsack problem parameters
	weights := []int{10, 20, 30, 40, 50}    // Weights of items
	values := []int{60, 100, 120, 140, 160} // Values of items
	capacity := 100                         // Knapsack capacity

	n := len(weights) // Number of items

	// Create binary variables for each item (0 if not selected, 1 if selected)
	vars := make([]*mip.Variable, n)
	for i := 0; i < n; i++ {
		vars[i] = solver.VarBool(fmt.Sprintf("x%d", i))
	}

	// Add the capacity constraint
	exp := mip.NewLinearExpression()
	for i := 0; i < n; i++ {
		exp.AddTerm(vars[i], float64(weights[i]))
	}
	solver.AddConstraintExpr(exp, mip.LessThanOrEqual, float64(capacity))

	// Set the objective function (maximize total value)
	obj := mip.NewLinearExpression()
	for i := 0; i < n; i++ {
		obj.AddTerm(vars[i], float64(values[i]))
	}
	solver.SetObjective(obj, mip.Maximize)

	// Solve the problem with no time limit
	foundOptimal, err := solver.Solve(-1)
	if err != nil {
		log.Fatalf("Error solving the problem: %v", err)
	}

	if !foundOptimal {
		fmt.Println("Solver finished within the time limit without finding the optimal solution.")
	} else {
		fmt.Println("Optimal solution found")
	}
}
