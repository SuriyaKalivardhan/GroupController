package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

func NewCorrector(redisClient *redis.Client, id string) *Corrector {
	corrector := &Corrector{redisClient, make(chan int), id}
	go corrector.correctnessWorker()
	return corrector
}

type Corrector struct {
	redisClient         *redis.Client
	currentCountChannel chan int
	id                  string
}

func (c *Corrector) correctnessWorker() {
	lastReconciliationTime := time.Now()
	for {
		select {
		case count := <-c.currentCountChannel:
			if time.Since(lastReconciliationTime) > (5 * time.Second) {
				lastReconciliationTime = time.Now()
				target, err := c.redisClient.Get("TARGET").Result()
				if err != nil {
					log.Printf("Unexpected error while fetching TARGET key.. %v", err)
					continue
				}
				targetInt, err := strconv.Atoi(target)
				if err != nil {
					log.Printf("Unexpected error while parsing TARGET key.. %v", err)
					continue
				}
				if targetInt != count {
					log.Printf("Request target %v while current count is %v", targetInt, count)
					SubmitTarget(c.id, targetInt)
				}
			}
		}
	}
}

func SubmitTarget(id string, target int) {
	request := AllocateRequest{
		Id:              id,
		Target:          target,
		RedisHost:       "localhost",
		RedisPort:       6388,
		RedisPassword:   "",
		RegisterChannel: ControllerBootChannel,
	}

	requestBytes, err := json.Marshal(request)

	if err != nil {
		log.Printf("Unexpected error while Serializing the request.. %v", err)
		return
	}

	resp, err := http.Post("http://localhost:5001/allocate", "application/json", bytes.NewBuffer(requestBytes))

	if err != nil {
		log.Printf("Unexpected error while Sending POST request.. %v", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unexpected error while reading Response body.. %v", err)
		return
	}
	log.Printf("RECEIVED response %v", string(body))
}
