package userstore

// UserStore is an exported store struct without interface assertion in store.go
type UserStore struct{} // want `LINT-013: store struct "UserStore" missing compile-time interface assertion in store\.go`
