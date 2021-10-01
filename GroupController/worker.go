package main

import (
	"log"
)

type worker struct {
	id              string
	listenerChannel string
	client          string
	close           func()
}

func NewWorker(message BootChannelMessage, close func()) *worker {
	return &worker{
		id:              message.Id,
		listenerChannel: message.ListenerChannel,
		client:          message.BindedClient,
		close:           close,
	}
}

func (w *worker) Update(message *BootChannelMessage) {
	if message.Id != w.id {
		log.Printf("Mismatch boot assignment %v for %v", message.Id, w.id)
		return
	}
	if message.ListenerChannel != w.listenerChannel {
		log.Printf("Mismatch boot assignment %v for %v", message.Id, w.id)
		w.listenerChannel = message.ListenerChannel
	}
	if message.BindedClient != w.client {
		log.Printf("Change in client from  %v to %v, reassigning", w.client, message.BindedClient)
		w.client = message.BindedClient
	}
}
