package mip

/*
#cgo CXXFLAGS: -std=c++17 -I/usr/local/include
#cgo LDFLAGS: -L${SRCDIR}/../bridge -lbridge -L/usr/local/lib -lortools -Wl,-rpath,/usr/local/lib
#include <stdlib.h>
#include "../bridge/bridge.h"
*/
import "C"
import (
	"unsafe"
)

// This file is a wrapper around the extern C functions defined in bridge.h.
// This is a direct translation of the bridge.h API to Go.
// Nothing from this file is exported outside of this package in order
// to separate the translation code from the actual exported mip API.

type solver struct{ csolver *C.CSolver }

func createSolver(solverType string) *solver {
	cName := C.CString(solverType)
	defer C.free(unsafe.Pointer(cName))
	return &solver{C.CreateSolver(cName)}
}

func (s *solver) delete()                     { C.DeleteSolver(s.csolver) }
func (s *solver) setMaximization()            { C.SetMaximization(s.csolver) }
func (s *solver) setMinimization()            { C.SetMinimization(s.csolver) }
func (s *solver) setTimeLimit(duration int64) { C.SetTimeLimit(s.csolver, C.int(duration)) }
func (s *solver) solve() int                  { return int(C.Solve(s.csolver)) }
func (s *solver) objectiveValue() float64     { return float64(C.ObjectiveValue(s.csolver)) }
func (s *solver) getBestBound() float64       { return float64(C.GetBestBound(s.csolver)) }
func (s *solver) setObjectiveCoefficient(variable *variable, coeff float64) {
	C.SetObjectiveCoefficient(s.csolver, variable.cvariable, C.double(coeff))
}

type variable struct{ cvariable *C.CVariable }

func (s *solver) newVariable(name string, lb, ub float64, varType int) *variable {
	// varType: 0 - continuous, 1 - integer
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return &variable{C.AddVar(s.csolver, cName, C.double(lb), C.double(ub), C.int(varType))}
}

func (v *variable) solutionValue() float64 { return float64(C.SolutionValue(v.cvariable)) }
func (v *variable) name() string           { return C.GoString(C.VariableName(unsafe.Pointer(v.cvariable))) }

type constraint struct{ cconstraint *C.CConstraint }

func (s *solver) newConstraint(lb, ub float64) *constraint {
	return &constraint{cconstraint: C.AddConstraint(s.csolver, C.double(lb), C.double(ub))}
}

func (c *constraint) setCoefficient(v *variable, coeff float64) {
	C.SetCoefficient(c.cconstraint, v.cvariable, C.double(coeff))
}
