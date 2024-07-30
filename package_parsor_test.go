package main

import "testing"

func TestParsePackage(t *testing.T) {
	// t.Skip("Not implemented")
	ParseTypes([]string{"UserRepository"}, []string{"repository/testdata"}, []string{})
}
