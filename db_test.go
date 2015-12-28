package main

import (
	"testing"
)

func TestSet(t *testing.T) {
	db := initDB("users_test.db")
	defer db.Close()
	err := db.set([]byte("test"))
	if err != nil {
		t.Errorf("set returned error: %s", err)
	}

	err = db.get([]byte("test"))
	if err != nil {
		t.Errorf("get returned error: %s", err)
	}

}
