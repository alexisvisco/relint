package entity

type Profile struct{}

type User struct {
	Profile Profile `gorm:"foreignKey:ProfileID"` // ok: only package model is checked
}
