package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/cmd/crud/handlers"
)

// monitoring
// https://github.com/rakyll/hey
// hey -m GET -c 10 -n 10000 "http://localhost:3000/v1/users"
//
// https://github.com/divan/expvarmon
// expvarmon -ports=":4000" -vars="requests,goroutines,errors,mem:memstats.Alloc"

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	// Configuration
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	shutdownTimeout := 5 * time.Second
	dbDialTimeout := 5 * time.Second

	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = ":3000"
	}

	debugHost := os.Getenv("DEBUG_HOST")
	if debugHost == "" {
		debugHost = ":4000"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		//set default dbhost
		dbHost = "localhost"
	}

	// Start mongodb
	log.Println("main started: Initialize Mongo")
	masterDB, err := db.New(dbHost, dbDialTimeout)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}
	defer masterDB.Close()

	// /debug/vars - Added to the default mux by the expvars package
	// /debug/pprof Added to the default mux by the net/http/pprof package
	debug := http.Server{
		Addr:           debugHost,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("startup : Debug Listening %s", debugHost)
		log.Printf("shutdown : Debug Listener closed : %v", debug.ListenAndServe())
	}()

	log.Println("main started: Initialize Server")
	// Start service
	api := http.Server{
		Addr:           apiHost,
		Handler:        handlers.API(masterDB),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// channel to collect errors from server
	// make it buffered so that go routine can exit even if we dont't collect the error
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("startup : Listening %s", apiHost)
		serverErrors <- api.ListenAndServe()
	}()

	// make a channel to listen to interupt or terminate signal from the OS
	// signal package requires a buffered channel
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting servevr: %v", err)
	case <-osSignals:
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown didnot complete in %v : %v", shutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Fatalf("Could not stop http server: %v", err)
			}
		}
	}

	log.Println("main completed")
}
