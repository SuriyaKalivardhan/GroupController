package main

import (
	"log"
	"time"
)

func NewCorrector() *Corrector {
	corrector := &Corrector{make(chan int)}
	go corrector.correctnessWorker()
	return corrector
}

type Corrector struct {
	currentCountChannel chan int
}

func (c *Corrector) correctnessWorker() {
	lastReconciliationTime := time.Now()
	for {
		select {
		case count := <-c.currentCountChannel:
			if time.Since(lastReconciliationTime) > (5 * time.Second) {
				log.Printf("Reconciling for %v", count)
				lastReconciliationTime = time.Now()
			}
		}
	}
}
