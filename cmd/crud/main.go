package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/cmd/crud/handlers"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	// Configuration
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	shutdownTimeout := 5 * time.Second
	host := os.Getenv("HOST")
	if host == "" {
		host = ":3000"
	}

	// Start service
	server := http.Server{
		Addr:           host,
		Handler:        handlers.API(),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Starting the service, listening for requests
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		log.Printf("startup : Listening %s", host)
		log.Printf("shutdown : Listener closed : %v", server.ListenAndServe())
		wg.Done()
	}()

	// Shutdown
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	<-osSignals

	// context for shutdown call
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Asking listenter to shutdown and load shed
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown didnot complete in %v : %v", shutdownTimeout, err)

		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error kiling server : %v", err)
		}
	}

	wg.Wait()
	log.Println("main completed")
}
