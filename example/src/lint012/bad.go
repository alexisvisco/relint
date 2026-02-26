package lint012

// This package simulates a store package.
// Since we can't import core/model in test data without a go.mod,
// this test just verifies the analyzer runs without errors in a store package.

type UserStore struct{}

func (s *UserStore) GetUser() string { // ok - not core/model type
	return ""
}
