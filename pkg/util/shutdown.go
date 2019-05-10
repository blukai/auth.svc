package util

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RegisterShutdown(timeout time.Duration, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	}

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, signals...)

	log.Printf("got signal: %+v. terminating...\n", <-quit)

	go func() {
		select {
		case <-time.After(timeout):
			log.Fatalln("timeout reached, terminating...")
		case s := <-quit:
			log.Fatalf("got signal: %v. terminating...\n", s)
		}
	}()
}
