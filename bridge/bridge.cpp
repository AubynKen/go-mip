#define BUILDING_BRIDGE
#include "bridge.h"
#include <ortools/linear_solver/linear_solver.h>

// This is a C interface to the OR-Tools linear solver. It is a simple wrapper around the C++ API
// without syntax sugar or error handling. The goal is to provide a minimalistic interface that can
// be used in Go with CGO

namespace {
    using Solver = operations_research::MPSolver;
    using Variable = operations_research::MPVariable;
    using Constraint = operations_research::MPConstraint;
}

extern "C" {
CSolver *CreateSolver(const char *solver_type) {
    Solver *solver = Solver::CreateSolver(solver_type);
    return reinterpret_cast<CSolver *>(solver);
}

void DeleteSolver(CSolver *solver) {
    delete reinterpret_cast<Solver *>(solver);
}

CVariable *AddVar(CSolver *solver, const char *name, double lb, double ub, int is_integer) {
    auto *s = reinterpret_cast<Solver *>(solver);
    Variable *var = s->MakeVar(lb, ub, is_integer, name);
    return reinterpret_cast<CVariable *>(var);
}

const char* VariableName(void* variable) {
    Variable* var = (Variable*)variable;
    return var->name().c_str();
}

CConstraint *AddConstraint(CSolver *solver, double lb, double ub) {
    auto *s = reinterpret_cast<Solver *>(solver);
    auto *constraint = s->MakeRowConstraint(lb, ub);
    return reinterpret_cast<CConstraint *>(constraint);
}

void SetCoefficient(CConstraint *constraint, CVariable *var, double coeff) {
    auto *c = reinterpret_cast<Constraint *>(constraint);
    auto *v = reinterpret_cast<Variable *>(var);
    c->SetCoefficient(v, coeff);
}

void SetObjectiveCoefficient(CSolver *solver, CVariable *var, double coeff) {
    auto *s = reinterpret_cast<Solver *>(solver);
    auto *v = reinterpret_cast<Variable *>(var);
    s->MutableObjective()->SetCoefficient(v, coeff);
}

void SetMaximization(CSolver *solver) {
    auto *s = reinterpret_cast<Solver *>(solver);
    s->MutableObjective()->SetMaximization();
}

void SetMinimization(CSolver *solver) {
    auto *s = reinterpret_cast<Solver *>(solver);
    s->MutableObjective()->SetMinimization();
}

void SetTimeLimit(CSolver *solver, int time_limit_nanoseconds) {
    auto *s = reinterpret_cast<Solver *>(solver);
    const absl::Duration time_limit = absl::Nanoseconds(time_limit_nanoseconds);
    s->SetTimeLimit(time_limit);
}

int Solve(CSolver *solver) {
    auto *s = reinterpret_cast<Solver *>(solver);
    return s->Solve();
}

double ObjectiveValue(CSolver *solver) {
    auto *s = reinterpret_cast<Solver *>(solver);
    return s->Objective().Value();
}

double SolutionValue(CVariable *var) {
    auto *v = reinterpret_cast<Variable *>(var);
    return v->solution_value();
}
}
