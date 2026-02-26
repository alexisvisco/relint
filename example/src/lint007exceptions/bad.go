package lint007exceptions

type Status string

const (
	Active        Status = "active"  // ok - configured exception for lint007exceptions.Status
	StatusPending Status = "pending" // ok
)

type Color int

const (
	Red        Color = 1 // want `LINT-007: const "Red" must be prefixed with type name "Color"`
	ColorGreen Color = 2 // ok
)
