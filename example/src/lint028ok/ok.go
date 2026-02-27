package model

type User struct {
	ID string `gorm:"column:id"`

	Email string `gorm:"column:email" db:"email"`

	status string
}
