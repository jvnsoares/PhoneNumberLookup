package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	shutdownTime = time.Second * 5
)

func registerRoutes() {
	http.HandleFunc(numberLookupEndpoint, numberLookupHandler)
}

// StartHTTPServer starts an HTTP server.
// the server servers the endpoints defined in setupRoutes()
// this is a blocking function
// this function listen for sigterm to shut down the HTTP server gracefully
func StartHTTPServer(port string) (stop func(), err error) {
	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	registerRoutes()
	errChan := make(chan error)

	go func() {
		log.Printf("Starting HTTP server on port %s", port)
		er := server.ListenAndServe()
		if er != nil {
			errChan <- er
		}
	}()

	// we wait 500 millisecond to be sure server.ListenAndServe will not return an error
	select {
	case err := <-errChan:
		return nil, err

	case <-time.NewTimer(time.Millisecond * 500).C:
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
		defer cancel()

		log.Printf("Shutting down http server")
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("error shutting down http server: %v", err)
		}
	}, nil
}
