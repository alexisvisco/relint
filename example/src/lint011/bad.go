package types

type UserService interface { // ok
	GetUser()
}

type UserStore interface { // ok
	FindUser()
}

type MyHelper interface { // want `LINT-011: interface "MyHelper" in types package must end with "Service" or "Store"`
	Help()
}
