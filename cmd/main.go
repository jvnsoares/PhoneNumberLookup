package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"phone-number-lookup/api"
)

const (
	portEnvKey  = "HTTP_PORT"
	defaultPort = "8008"
)

func main() {
	port := os.Getenv(portEnvKey)
	if port == "" {
		port = defaultPort
	}

	stop, err := api.StartHTTPServer(port)
	if err != nil {
		log.Println(err)
		return
	}
	defer stop()

	// wait for sigterm to gracefully stop the http server
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
	}

	log.Println("sigterm received")
}
