package lint007

type Status string

const (
	Active          Status = "active"   // want `LINT-007: const "Active" must be prefixed with type name "Status"`
	StatusInactive  Status = "inactive" // ok
	StatusActive    Status = "active2"  // ok
	Pending         Status = "pending"  // want `LINT-007: const "Pending" must be prefixed with type name "Status"`
)

type Color int

const (
	Red         Color = 1 // want `LINT-007: const "Red" must be prefixed with type name "Color"`
	ColorGreen  Color = 2 // ok
	ColorBlue   Color = 3 // ok
)
