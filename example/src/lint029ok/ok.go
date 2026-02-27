package model

type (
	Profile struct{}
	Role    struct{}
)

type User struct {
	Profile *Profile `gorm:"foreignKey:ProfileID"`

	Roles []*Role `gorm:"many2many:user_roles"`

	Label string `gorm:"column:label"`
}
