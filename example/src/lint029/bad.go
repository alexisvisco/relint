package model

type (
	Profile struct{}
	Role    struct{}
	Tag     struct{}
)

type User struct {
	Profile Profile `gorm:"foreignKey:ProfileID"` // want `LINT-029: relation field "Profile" with gorm tag "foreignKey" must be a pointer or a slice of pointers \(\[\]\*Type\)`

	Roles []Role `gorm:"many2many:user_roles"` // want `LINT-029: relation field "Roles" with gorm tag "many2many" must be a pointer or a slice of pointers \(\[\]\*Type\)`

	Tag Tag `gorm:"polymorphicType:OwnerType"` // want `LINT-029: relation field "Tag" with gorm tag "polymorphicType" must be a pointer or a slice of pointers \(\[\]\*Type\)`

	Manager *Profile `gorm:"foreignKey:ManagerID"`

	RoleRefs []*Role `gorm:"many2many:user_role_refs"`

	Label string `gorm:"column:label"`
}
