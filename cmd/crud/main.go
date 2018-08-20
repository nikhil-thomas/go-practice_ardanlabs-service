package main

import (
	"context"
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/cfg"

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

	c, err := cfg.New(cfg.EnvProvider{Namespace: "CRUD"})
	if err != nil {
		log.Printf("%s all configs defaults in use", err)
	}

	readTimeout, err := c.Duration("READ_TIMEOUT")
	if err != nil {
		readTimeout = 5 * time.Second
	}

	writeTimeout, err := c.Duration("WRITE_TIMEOUT")
	if err != nil {
		writeTimeout = 5 * time.Second
	}

	shutdownTimeout, err := c.Duration("SHUTDOWN_TIMEOUT")
	if err != nil {
		shutdownTimeout = 5 * time.Second
	}

	dbDialTimeout, err := c.Duration("DB_DIAL_TIMEOUT")
	if err != nil {
		dbDialTimeout = 5 * time.Second
	}

	apiHost, err := c.String("API_HOST")
	if err != nil {
		apiHost = ":3000"
	}

	debugHost, err := c.String("DEBUG_HOST")
	if err != nil {
		debugHost = ":4000"
	}

	dbHost, err := c.String("DB_HOST")
	if dbHost == "" {
		//set default dbhost
		dbHost = "localhost"
	}
	log.Printf("%s=%v", "READ_TIMEOUT", readTimeout)
	log.Printf("%s=%v", "WRITE_TIMEOUT", writeTimeout)
	log.Printf("%s=%v", "SHUTDOWN_TIMEOUT", shutdownTimeout)
	log.Printf("%s=%v", "DB_DIAL_TIMEOUT", dbDialTimeout)
	log.Printf("%s=%v", "API_HOST", apiHost)
	log.Printf("%s=%v", "DEBUG_HOST", debugHost)
	log.Printf("%s=%v", "DB_HOST", dbHost)

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

	log.Println("main : Start shutdown...")

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
