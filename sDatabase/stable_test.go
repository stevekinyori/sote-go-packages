package sDatabase

import (
	"testing"
)

func TestGetTables(t *testing.T) {
	soteErr := GetConnection("single", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Should have returned a pointer to the single database connection")
		t.Fail()
	}

	getTables("sote")

	soteErr = GetConnection("pool", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Should have returned a pointer to the single database connection")
		t.Fail()
	}

	getTables("sote")
}
