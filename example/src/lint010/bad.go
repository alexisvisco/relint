package lint010

type UserService interface { // want `LINT-010: interface "UserService" must be declared in a types package`
	DoSomething()
}

type UserStore interface { // want `LINT-010: interface "UserStore" must be declared in a types package`
	GetByID()
}

type DetailErrorer interface{} // ok - not a Service/Store interface

type ConcreteStruct struct{}
