package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

func NewWorkerInteractor(client *redis.Client, rootContext context.Context) (*WorkerInteractor, context.CancelFunc) {
	controllerContext, cancelFunc := context.WithCancel(rootContext)
	w := &WorkerInteractor{client, controllerContext, make(map[string]*worker), NewCorrector(client)}
	go w.workerWatcher()
	return w, cancelFunc
}

type WorkerInteractor struct {
	redisClient *redis.Client
	ctx         context.Context
	workers     map[string]*worker
	corrector   *Corrector
}

func (w *WorkerInteractor) handleRedisMessage(redisPayload string) {
	var message BootChannelMessage
	err := json.Unmarshal([]byte(redisPayload), &message)
	if err != nil {
		log.Printf("Unexpected REDIS error message %v", err)
		return
	}

	if w.workers[message.Id] != nil {
		if !w.workers[message.Id].Update(&message) {
			w.workers[message.Id].close()
			delete(w.workers, message.Id)
			log.Printf("Shutting down %v on receive from worker", message.Id)
		}
	} else {
		log.Printf("Adding new worker %v", message.Id)
		_, cancel := context.WithCancel(w.ctx)
		worker := NewWorker(message, cancel)
		w.workers[worker.id] = worker
	}
}

func (w *WorkerInteractor) workerWatcher() {
	for {
		select {
		case <-w.ctx.Done():
			log.Println("context closed for Controller in WatchWorkerListern, returning")
			return
		case <-time.After(1 * time.Second):
			var channels []string
			localwokers := make(map[string]*worker, len(w.workers))
			for _, worker := range w.workers {
				channels = append(channels, worker.listenerChannel)
				localwokers[worker.listenerChannel] = worker
			}

			currentCount := 0
			if len(channels) == 0 {
				w.corrector.currentCountChannel <- currentCount
				continue
			}
			result, err := w.redisClient.PubSubNumSub(channels...).Result()
			if err != nil {
				log.Printf("Unexpected error while Wathicng on listeners %v", err)
				continue
			}

			for channel, count := range result {
				if count > 0 {
					currentCount++
				}

				if localwokers[channel] == nil {
					continue
				}

				if count == 0 {
					log.Printf("No subsribre for %v, removing %s", channel, localwokers[channel].id)
					delete(w.workers, localwokers[channel].id)
					localwokers[channel].close()
				}
			}

			w.corrector.currentCountChannel <- currentCount
		}
	}
}
