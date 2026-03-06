package userstore

type UserStore struct{}

func New() *UserStore { // want `LINT-032: package "userstore" must expose only one constructor matching New\*; found 2`
	return &UserStore{}
}

func NewUserStore() *UserStore { // want `LINT-032: constructor "NewUserStore" in package "userstore" must be named "New"`
	return &UserStore{}
}
