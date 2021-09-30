package main

import "testing"

func TestHello(t *testing.T) {
	result := getHello("Agazhi")
	if result != "Hello Agazhi" {
		t.Fatalf("Unexepcted result %s", result)
	}
}

func TestSampleRequest(t *testing.T) {
	request := getSampleRquest()
	if request.id != "suriya" && request.target != 34 {
		t.Fatalf("Unexpected result %s, %v", request.id, request.target)
	}

	nextRequest := copyRequest(request)
	nextRequest.id = "magi"
	nextRequest.target = 31

	if request.id != "suriya" || request.target != 34 {
		t.Fatalf("Unexpected result %s, %v", request.id, request.target)
	}

	if nextRequest.id != "magi" || nextRequest.target != 31 {
		t.Fatalf("Unexpected result %s, %v", nextRequest.id, nextRequest.target)
	}

	pointerRequest := request
	if request.id != "suriya" || request.target != 34 {
		t.Fatalf("Unexpected result %s, %v", request.id, request.target)
	}

	if pointerRequest.id != "suriya" || pointerRequest.target != 34 {
		t.Fatalf("Unexpected result %s, %v", pointerRequest.id, pointerRequest.target)
	}

	pointerRequest.id = "agazhi"
	pointerRequest.target = 4
	if request.id != "agazhi" || request.target != 4 {
		t.Fatalf("Unexpected result %s, %v", request.id, request.target)
	}

	if pointerRequest.id != "agazhi" || pointerRequest.target != 4 {
		t.Fatalf("Unexpected result %s, %v", pointerRequest.id, pointerRequest.target)
	}

}
