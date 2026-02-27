package entity

type User struct {
	ID string `json:"id"` // ok: only package model is checked
}
