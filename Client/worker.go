package main

import (
	"log"
)

type worker struct {
	id              string
	listenerChannel string
	close           func()
}

func NewWorker(message BootChannelMessage, close func()) *worker {
	return &worker{
		id:              message.Id,
		listenerChannel: message.ListenerChannel,
		close:           close,
	}
}

func (w *worker) Update(message *BootChannelMessage) bool {
	if message.Method == "Shutdown" {
		return false
	}

	if message.Id != w.id {
		log.Printf("Mismatch boot assignment %v for %v", message.Id, w.id)
		return true
	}

	if message.ListenerChannel != w.listenerChannel {
		log.Printf("Mismatch boot assignment %v for %v", message.Id, w.id)
		w.listenerChannel = message.ListenerChannel
	}

	return true
}
