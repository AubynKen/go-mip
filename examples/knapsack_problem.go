package examples

import (
	"fmt"
	"gomip/mip"
	"log"
)

func KnapsackProblem() {
	// Create a new solver
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		log.Fatalf("Error creating solver: %v", err)
	}
	defer solver.ReleaseResources()

	// Define the knapsack problem parameters
	weights := []int{10, 20, 30, 40, 50, 25, 1}      // Weights of items
	values := []int{60, 100, 120, 140, 160, 130, 10} // Values of items
	capacity := 100                                  // Knapsack capacity

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

	// Print out the solution
	totalValue := 0
	totalWeight := 0
	selected := make([]int, 0)
	for i := 0; i < n; i++ {
		if vars[i].Value() > 0.5 {
			selected = append(selected, i)
			totalValue += values[i]
			totalWeight += weights[i]
		}
	}

	fmt.Printf("Selected items: %v\n", selected)
	fmt.Printf("Total value: %d\n", totalValue)
	fmt.Printf("Total weight: %d\n", totalWeight)
}
