package model

type User struct {
	ID string // want `LINT-028: exported model field "ID" must declare a gorm tag`

	Email string `db:"email"` // want `LINT-028: exported model field "Email" must declare a gorm tag`

	Name string `gorm:"column:name"`

	status string
}
