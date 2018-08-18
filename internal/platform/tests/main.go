package tests

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"testing"
	"time"

	"github.com/pborman/uuid"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/docker"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
)

// Success and failure markers
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// Test owns state for running/shutting down tests
type Test struct {
	MasterDB  *db.DB
	container *docker.Container
}

// New is the entry poinbt for tests
func New() *Test {
	var test Test

	var err error

	test.container, err = docker.StartMongo()

	if err != nil {
		log.Fatalln(err)
	}

	dbDialTimeout := 25 * time.Second
	dbHost := fmt.Sprintf("mongodb://localhost:%s/gotraining", test.container.Port)

	log.Println("main: started: Initialize Mongo")
	test.MasterDB, err = db.New(dbHost, dbDialTimeout)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}

	return &test
}

// TearDown shutting down tests
func (t *Test) TearDown() {
	t.MasterDB.Close()
	if err := docker.StopMongo(t.container); err != nil {
		log.Println(err)
	}
}

// Recover prevensts panics from blocking the test from cleanup
func Recover(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal("Unhandled Exception:", string(debug.Stack()))
	}
}

// Context returns an app level context for testing
func Context() context.Context {
	values := web.Values{
		TraceID: uuid.New(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}
