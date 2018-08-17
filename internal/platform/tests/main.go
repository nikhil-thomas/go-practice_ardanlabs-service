package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/docker"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"

	"github.com/pborman/uuid"
)

// Test set up
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

// MasterDB represents the master MOngo session
// This will work since we are starting mongo from a container
var MasterDB *db.DB

// Context returns returns an app level context for testing
func Context() context.Context {
	values := web.Values{
		TraceID: uuid.New(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}

// Main is the entry poinbt for tests
func Main(m *testing.M) int {

	c, err := docker.StartMongo()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := docker.StopMongo(c); err != nil {
			log.Println(err)
		}
	}()

	dbDialTimeout := 25 * time.Second
	dbHost := fmt.Sprintf("mongodb://localhost:%s/gotraining", c.Port)

	log.Println("main: started: Initialize Mongo")
	MasterDB, err = db.New(dbHost, dbDialTimeout)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}

	defer MasterDB.Close()

	return m.Run()
}
