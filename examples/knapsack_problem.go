package examples

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"gomip/mip"
	"log"
	"os"
)

func KnapsackProblem() {
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		log.Fatalf("Error creating solver: %v", err)
	}
	defer solver.ReleaseResources()

	weights := []int{10, 20, 30, 40, 50, 25, 1}      // Weights of items
	values := []int{60, 100, 120, 140, 160, 130, 10} // Values of items
	capacity := 100                                  // Knapsack capacity
	n := len(weights)                                // Number of items

	// var[i] = 1 if item i is selected, 0 otherwise
	vars := make([]*mip.Variable, n)
	for i := 0; i < n; i++ {
		vars[i] = solver.VarBool(fmt.Sprintf("x%d", i))
	}

	// total weight should be less than knapsack capacity
	exp := mip.NewLinearExpression()
	for i := 0; i < n; i++ {
		exp.AddTerm(vars[i], float64(weights[i]))
	}
	solver.AddConstraintExpr(exp, mip.LessThanOrEqual, float64(capacity))

	// objective: maximize total value
	obj := mip.NewLinearExpression()
	for i := 0; i < n; i++ {
		obj.AddTerm(vars[i], float64(values[i]))
	}
	solver.SetObjective(obj, mip.Maximize)

	foundOptimal, err := solver.Solve(-1) // no time limit
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

	knapsackProblemOutput(weights, values, selected, capacity)
}

// knapsackProblemOutput prints the results of the knapsack problem in a tabular format.
// No need to dig into the details of this function.
func knapsackProblemOutput(weights []int, values []int, selected []int, capacity int) {
	t := table.NewWriter()
	t.SetCaption("Knapsack Problem Results")
	t.SetOutputMirror(os.Stdout)

	// Add header
	header := table.Row{"Object"}
	for i := range weights {
		header = append(header, fmt.Sprintf("#%d", i))
	}
	header = append(header, "Total")
	t.AppendHeader(header)

	// Add weight row
	weightRow := table.Row{"Weight"}
	totalWeight := 0
	for _, w := range weights {
		weightRow = append(weightRow, w)
		totalWeight += w
	}
	weightRow = append(weightRow, totalWeight)
	t.AppendRow(weightRow)

	// Add value row
	valueRow := table.Row{"Value"}
	totalValue := 0
	for _, v := range values {
		valueRow = append(valueRow, v)
		totalValue += v
	}
	valueRow = append(valueRow, totalValue)
	t.AppendRow(valueRow)
	t.Render()

	fmt.Println() // Add a newline for better readability

	// Create table for results
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Metric", "Value"})

	selectedWeight := 0
	selectedValue := 0
	for _, i := range selected {
		selectedWeight += weights[i]
		selectedValue += values[i]
	}

	t.AppendRows([]table.Row{
		{"Selected items", selected},
		{"Total selected weight", selectedWeight},
		{"Total selected value", selectedValue},
		{"Knapsack capacity", capacity},
	})

	t.Render()
}
