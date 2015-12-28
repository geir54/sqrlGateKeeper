package main

import (
	"testing"
)

func TestNuts(t *testing.T) {
	n := initnutList(2)
	n.add("test1")
	n.add("test2")

	data, ok := n.get("test2")
	if !ok {
		t.Fatal("test2 should be set")
	}

	data.Autenticated = true
	n.update("test2", data)

	data2, ok := n.get("test2")
	if !ok || !data2.Autenticated {
		t.Fatal("test2 should be Autenticated")
	}
}
