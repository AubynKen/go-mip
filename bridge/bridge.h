#ifndef BRIDGE_H
#define BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

#if defined(_WIN32) || defined(__CYGWIN__)
  #ifdef BUILDING_BRIDGE
    #define BRIDGE_API __declspec(dllexport)
  #else
    #define BRIDGE_API __declspec(dllimport)
  #endif
#else
  #define BRIDGE_API __attribute__((visibility("default")))
#endif

typedef void* CSolver;
typedef void* CVariable;
typedef void* CConstraint;

BRIDGE_API CSolver *CreateSolver(const char *solver_type);
BRIDGE_API void DeleteSolver(CSolver* solver);
BRIDGE_API CVariable* AddVar(CSolver* solver, const char* name, double lb, double ub, int is_integer);
BRIDGE_API CConstraint* AddConstraint(CSolver* solver, double lb, double ub);
BRIDGE_API void SetCoefficient(CConstraint* constraint, CVariable* var, double coeff);
BRIDGE_API void SetObjectiveCoefficient(CSolver* solver, CVariable* var, double coeff);
BRIDGE_API const char* VariableName(void* variable);
BRIDGE_API void SetMaximization(CSolver* solver);
BRIDGE_API void SetMinimization(CSolver* solver);
BRIDGE_API void SetTimeLimit(CSolver *solver, int time_limit_nanoseconds);
BRIDGE_API int Solve(CSolver* solver);
BRIDGE_API double ObjectiveValue(CSolver* solver);
BRIDGE_API double SolutionValue(CVariable* var);

#ifdef __cplusplus
}
#endif

#endif // BRIDGE_H
