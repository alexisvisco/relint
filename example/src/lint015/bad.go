package userstore // want `LINT-015: file "bad.go" in store/service/handler package must contain exactly one exported store/service/handler method, found 2`

type UserStore struct{}

func (s *UserStore) CreateUser() {}

func (s *UserStore) UpdateUser() {}

func BuildPasswordHash() {} // ok - exported function but not a store/service method
