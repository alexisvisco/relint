package model

type User struct {
	ID string `json:"id"` // want `LINT-027: model struct fields must not declare json tags`

	Email string `json:"email" db:"email"` // want `LINT-027: model struct fields must not declare json tags`

	Age int `db:"age" json:"age,omitempty" validate:"required"` // want `LINT-027: model struct fields must not declare json tags`

	Status string `db:"status"`
}
