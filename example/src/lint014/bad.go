package userservice

// UserService is an exported service struct without interface assertion in service.go
type UserService struct{} // want `LINT-014: service struct "UserService" missing compile-time interface assertion in service\.go`
