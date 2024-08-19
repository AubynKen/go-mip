package examples

import (
	"fmt"
	"math"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"gomip/mip"
)

func TransportationProblem() {
	problem := createTransportationProblemData()
	printProblemData(problem)
	solution, totalCost := solveTransportationProblem(problem)
	printSolution(problem, solution, totalCost)
}

type transportationProblemData struct {
	Sources      []string
	Destinations []string
	Supply       map[string]float64
	Demand       map[string]float64
	Cost         map[string]map[string]float64
}

func solveTransportationProblem(problem transportationProblemData) (map[string]map[string]float64, float64) {
	// Create a new solver
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		fmt.Printf("Error creating solver: %v\n", err)
		return nil, 0
	}
	defer solver.ReleaseResources()

	// Create variables for transportation quantities
	vars := make(map[string]map[string]*mip.Variable)
	for _, source := range problem.Sources {
		vars[source] = make(map[string]*mip.Variable)
		for _, dest := range problem.Destinations {
			vars[source][dest] = solver.VarFloat(fmt.Sprintf("x_%s_%s", source, dest), 0, math.MaxFloat64)
		}
	}

	// supply constraints
	for _, source := range problem.Sources {
		exp := mip.NewLinearExpression()
		for _, dest := range problem.Destinations {
			exp.AddVar(vars[source][dest])
		}
		solver.AddConstraintExpr(exp, mip.LessThanOrEqual, problem.Supply[source])
	}

	// demand constraints
	for _, dest := range problem.Destinations {
		exp := mip.NewLinearExpression()
		for _, source := range problem.Sources {
			exp.AddVar(vars[source][dest])
		}
		solver.AddConstraintExpr(exp, mip.GreaterThanOrEqual, problem.Demand[dest])
	}

	// minimize total transportation cost
	obj := mip.NewLinearExpression()
	for _, source := range problem.Sources {
		for _, dest := range problem.Destinations {
			obj.AddTerm(vars[source][dest], problem.Cost[source][dest])
		}
	}
	solver.SetObjective(obj, mip.Minimize)

	_, err = solver.Solve(-1)
	if err != nil {
		fmt.Printf("Error solving the problem: %v\n", err)
		return nil, 0
	}

	// Extract solution
	solution := make(map[string]map[string]float64)
	for _, source := range problem.Sources {
		solution[source] = make(map[string]float64)
		for _, dest := range problem.Destinations {
			solution[source][dest] = vars[source][dest].Value()
		}
	}

	return solution, solver.ObjectiveValue()
}

func createTransportationProblemData() transportationProblemData {
	return transportationProblemData{
		Sources:      []string{"Factory1", "Factory2"},
		Destinations: []string{"Store1", "Store2", "Store3"},
		Supply: map[string]float64{
			"Factory1": 100,
			"Factory2": 150,
		},
		Demand: map[string]float64{
			"Store1": 80,
			"Store2": 70,
			"Store3": 90,
		},
		Cost: map[string]map[string]float64{
			"Factory1": {"Store1": 2, "Store2": 3, "Store3": 1},
			"Factory2": {"Store1": 5, "Store2": 4, "Store3": 6},
		},
	}
}

// for pretty printing only, no need to dig into the details of the function
func printProblemData(problem transportationProblemData) {
	fmt.Println("Transportation Problem Data, supply and demand of factories / stores:")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "Supply/Demand"})

	for _, source := range problem.Sources {
		t.AppendRow(table.Row{source, problem.Supply[source]})
	}
	for _, dest := range problem.Destinations {
		t.AppendRow(table.Row{dest, problem.Demand[dest]})
	}
	t.Render()

	fmt.Println("\n\nTransportation Problem Data, cost matrix:")

	// Print cost matrix
	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{"From \\ To"}
	for _, dest := range problem.Destinations {
		header = append(header, dest)
	}
	t.AppendHeader(header)

	for _, source := range problem.Sources {
		row := table.Row{source}
		for _, dest := range problem.Destinations {
			row = append(row, problem.Cost[source][dest])
		}
		t.AppendRow(row)
	}

	t.Render()
}

// printSolution prints the solution in a pretty tabular format, no need to dig into the details of the function
func printSolution(problem transportationProblemData, solution map[string]map[string]float64, totalCost float64) {
	t := table.NewWriter()
	fmt.Printf(fmt.Sprintf("\n\nOptimal Transportation Plan (Total Cost: %.2f)\n", totalCost))
	t.SetOutputMirror(os.Stdout)

	header := table.Row{"From \\ To"}
	for _, dest := range problem.Destinations {
		header = append(header, dest)
	}
	t.AppendHeader(header)

	for _, source := range problem.Sources {
		row := table.Row{source}
		for _, dest := range problem.Destinations {
			row = append(row, fmt.Sprintf("%.2f", solution[source][dest]))
		}
		t.AppendRow(row)
	}

	t.Render()
}
