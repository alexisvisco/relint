package userservice

type UserService struct{}

func NewUserService() *UserService { // want `LINT-032: constructor "NewUserService" in package "userservice" must be named "New"`
	return &UserService{}
}
