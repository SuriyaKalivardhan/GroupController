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
	if request.Id != "suriya" && request.Target != 34 {
		t.Fatalf("Unexpected result %s, %v", request.Id, request.Target)
	}

	nextRequest := copyRequest(request)
	nextRequest.Id = "magi"
	nextRequest.Target = 31

	if request.Id != "suriya" || request.Target != 34 {
		t.Fatalf("Unexpected result %s, %v", request.Id, request.Target)
	}

	if nextRequest.Id != "magi" || nextRequest.Target != 31 {
		t.Fatalf("Unexpected result %s, %v", nextRequest.Id, nextRequest.Target)
	}

	pointerRequest := request
	if request.Id != "suriya" || request.Target != 34 {
		t.Fatalf("Unexpected result %s, %v", request.Id, request.Target)
	}

	if pointerRequest.Id != "suriya" || pointerRequest.Target != 34 {
		t.Fatalf("Unexpected result %s, %v", pointerRequest.Id, pointerRequest.Target)
	}

	pointerRequest.Id = "agazhi"
	pointerRequest.Target = 4
	if request.Id != "agazhi" || request.Target != 4 {
		t.Fatalf("Unexpected result %s, %v", request.Id, request.Target)
	}

	if pointerRequest.Id != "agazhi" || pointerRequest.Target != 4 {
		t.Fatalf("Unexpected result %s, %v", pointerRequest.Id, pointerRequest.Target)
	}

}
