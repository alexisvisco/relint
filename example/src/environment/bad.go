package environment

type Environment string

const (
	Development           Environment = "development" // ok - default exception for environment.Environment
	EnvironmentProduction Environment = "production"  // ok
)

type Status string

const (
	Active       Status = "active"  // want `LINT-007: const "Active" must be prefixed with type name "Status"`
	StatusPaused Status = "paused"  // ok
	Pending      Status = "pending" // want `LINT-007: const "Pending" must be prefixed with type name "Status"`
)
