package main

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func Test_GetRedisKey(t *testing.T) {
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Error initiating %v", err)
	}
	redisServer.Set("Agazhi", "Magi")
	client := initRedis(redisServer.Addr(), "")
	result := getValue(client, "Agazhi")
	if result != "Magi" {
		t.Fatalf("Unexpected result %v", result)
	}
}

func Test_GetRedisSubscribe(t *testing.T) {
	redisServer, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Error initiating %v", err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		redisServer.Publish("BootChannel", "Liveness message")
	}()

	client := initRedis(redisServer.Addr(), "")
	pubSub := client.Subscribe("BootChannel")
	select {
	case message := <-pubSub.Channel():
		if message.Payload != "Liveness message" {
			t.Fatalf("Error %v", message)
		}
	}
}
