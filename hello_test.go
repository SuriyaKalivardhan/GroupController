package main

import "testing"

func TestHello(t *testing.T) {
	result := getHello("Agazhi")
	if result != "Hello Agazhi" {
		t.Fatalf("Unexepcted result %s", result)
	}
}
