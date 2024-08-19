package examples

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"gomip/mip"
	"log"
	"math/rand"
	"os"
	"time"
)

func RoutingCtrlLinkSelection(timeLimit time.Duration) {
	// Set random seed for reproducibility
	rd := rand.New(rand.NewSource(42))

	// Define the problem parameters
	numGroups := 10  // number of prefix groups to fit links into
	numLinks := 600  // number of links, each associated with a single device
	numDevices := 40 // number of devices

	// Generate random capacities for links (in Mbps)
	capacities := make([]float64, numLinks)
	for i := range capacities {
		capacities[i] = float64(rd.Intn(800)+100) + rd.Float64()
	}

	// Generate wanted capacities for prefix groups
	// Let's choose 1000, 2000, 3000, ... for the wanted capacities
	wantedCapacities := make([]float64, numGroups)
	for i := range wantedCapacities {
		wantedCapacities[i] = 1000 * float64(i+1)
	}

	// Generate random latencies for links (in ms)
	latencies := make([]float64, numLinks)
	for i := range latencies {
		latencies[i] = rd.Float64()*49 + 1
	}

	// Generate random packet loss rates for links (in percentage)
	loss := make([]float64, numLinks)
	for i := range loss {
		loss[i] = rd.Float64()
	}

	// Generate random device capacities (in Mbps)
	C := make([]float64, numDevices)
	for i := range C {
		C[i] = float64(rd.Intn(4000) + 2000)
	}

	// Assign links to devices randomly
	S := make([]int, numLinks)
	for link := range S {
		S[link] = rd.Intn(numDevices)
	}

	deviceToLinks := make(map[int][]int)
	for link, device := range S {
		deviceToLinks[device] = append(deviceToLinks[device], link)
	}

	// CBC solver seems to perform fine for this problem
	// It is slower than Gurobi/CPLEX, but it is free and open-source with no license restrictions
	solver, err := mip.NewSolver(mip.CBC)
	if err != nil {
		panic(err)
	}
	defer solver.ReleaseResources()

	// selection[gr][lk] = 1 if and only if prefix group `gr` selects link `lk`
	selections := make([][]*mip.Variable, numGroups)
	for group := range selections {
		selections[group] = make([]*mip.Variable, numLinks)
		for link := range selections[group] {
			varName := fmt.Sprintf("link_%d_selected_by_group_%d", link, group)
			selections[group][link] = solver.VarBool(fmt.Sprintf(varName, group, link))
		}
	}

	// constraint: a link is either unused or selected by at most one group
	for link := 0; link < numLinks; link++ {
		selectionCount := mip.NewLinearExpression() // how many groups select this link
		for group := 0; group < numGroups; group++ {
			selectionCount.AddVar(selections[group][link])
		}
		solver.AddConstraintExpr(selectionCount, "<=", 1)
	}

	// usageUB[d] represents an upper bound for the usage fraction of device d
	// this is a helper variable that is at least as big as the usage fraction of the device when solution valid
	usageUB := make([]*mip.Variable, numDevices)
	for device := range usageUB {
		varName := fmt.Sprintf("usage_fraction_upper_bound_device_%d", device)
		usageUB[device] = solver.VarFloat(varName, 0, 1)

		// sum of all capacities of all links belonging to the device, whether selected or not
		totalCapacity := 0.
		for _, link := range deviceToLinks[device] {
			totalCapacity += capacities[link]
		}

		// used capacity = sum{(selected * capacity) for all links belonging to the device}
		usedCapacity := mip.NewLinearExpression()
		for _, link := range deviceToLinks[device] {
			for group := range selections {
				usedCapacity.AddTerm(selections[group][link], capacities[link])
			}
		}

		// Constraint: (used capacity)/(total capacity) <= (device usage upperbound)
		// i.e. {used capacity} - {total capacity} * {device usage upperbound} <= 0:
		leftHandSide := mip.NewLinearExpression()             // 0
		leftHandSide.AddExpr(usedCapacity)                    // +{used capacity}
		leftHandSide.AddTerm(usageUB[device], -totalCapacity) // -{total capacity} * {device usage upperbound}
		solver.AddConstraintExpr(leftHandSide, mip.LessThanOrEqual, 0)
	}

	// Variable: Global usage upperbound, to be used in the objective function
	// this will serve as max(usage fraction of all devices)
	globalUsageUB := solver.VarFloat("global_device_usage_upperbound", 0, 1)

	for _, deviceUsageUB := range usageUB {
		// local device usage upperbound <= global usage upperbound
		// (globalUB - deviceUB) <= 0
		ubDiff := mip.NewLinearExpression()
		ubDiff.AddTerm(deviceUsageUB, -1)
		ubDiff.AddVar(globalUsageUB)
		solver.AddConstraintExpr(ubDiff, mip.GreaterThanOrEqual, 0)
	}

	// Define weights / importance for the objective function
	alpha := 1.       // Weight for latency penalty
	beta := 100.      // Weight for packet loss penalty
	gamma := 1000000. // Weight for device usage penalty

	// Objective function: sum of performance penalties, weighted by their capacities, plus the max usage penalty
	// - links that are not selected do not affect our scores, no matter how bad their performance is
	// - links that are selected should influence the score based on their performance, and proportionally to their capacity
	// - the global usage upperbound should be minimized, to ensure that the usage is balanced across devices
	objective := mip.NewLinearExpression()
	for gr := range selections {
		for lk := range selections[gr] {
			performancePenalty := alpha*latencies[lk] + beta*loss[lk]
			objective.AddTerm(selections[gr][lk], performancePenalty*capacities[lk])
		}
	}
	objective.AddTerm(globalUsageUB, gamma)

	// We only have penalties in the objective function (the greater they are the worse they are), so we want to minimize the objective
	solver.SetObjective(objective, mip.Minimize)

	// Constraint: Capacity constraint for each prefix group
	// The total capacity of the selected links for each group should be at least the wanted capacity
	for i := range selections {
		expr := mip.NewLinearExpression()
		for j := range selections[i] {
			expr.AddTerm(selections[i][j], capacities[j])
		}
		solver.AddConstraintExpr(expr, mip.GreaterThanOrEqual, wantedCapacities[i])
	}

	// Solve the problem with a time limit
	// Note that the solver will return the best solution found within the time
	// limit, with or without optimality guarantees.
	isOptimal, err := solver.Solve(timeLimit)
	if err != nil {
		log.Fatalf("Solver error: %v", err)
	}

	if isOptimal {
		fmt.Println("The objective is guaranteed to be optimal.")
	} else {
		fmt.Println("Suboptimal feasible solution found within time limit.")
	}

	// Print the results
	summary := table.NewWriter()
	summary.SetOutputMirror(os.Stdout)
	summary.AppendHeader(table.Row{
		"Best Objective Found",
		"Max Device Usage Fraction",
		"Lower Bound",
		"Gap (%)",
	})

	gapPercentage := solver.Gap() * 100
	summary.AppendRow(table.Row{
		fmt.Sprintf("%.2f", solver.ObjectiveValue()),
		fmt.Sprintf("%.2f", globalUsageUB.Value()),
		fmt.Sprintf("%.2f", solver.BestBound()),
		fmt.Sprintf("%.2f%%", gapPercentage),
	})
	summary.Render()

	fmt.Printf("\nHaving a gap of %.2f%% means that the objective value is proven to be at most within at most %.2f%% of the optimal.\n", gapPercentage, gapPercentage)

	fmt.Println("\nDetailed information about the link selections:")
	prefixGroupTable := table.NewWriter()
	prefixGroupTable.SetOutputMirror(os.Stdout)
	prefixGroupTable.AppendHeader(table.Row{
		"Prefix Group",
		"Wanted Capacity",
		"Actual Capacity",
		"Selected Links",
	})

	for gr, row := range selections {
		selected := make([]int, 0, len(row))
		totalCapacity := 0.0

		for lk, varLink := range row {
			if varLink.Value() > 0.5 {
				selected = append(selected, lk)
				totalCapacity += capacities[lk]
			}
		}

		selectedLinksStr := fmt.Sprintf("%v", selected)
		prefixGroupTable.AppendRow(table.Row{
			gr,
			fmt.Sprintf("%.2f", wantedCapacities[gr]),
			fmt.Sprintf("%.2f", totalCapacity),
			selectedLinksStr,
		})
	}
	prefixGroupTable.Render()
}
