package user_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/docker"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/web"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/user"
	"github.com/pborman/uuid"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
)

func TestUser(t *testing.T) {
	t.Log("Given the need to validate creating a user")
	{
		t.Log("\t when handling a single user")
		{
			dbConn, err := masterDB.Copy()
			if err != nil {
				t.Fatalf("\t%s\tshould be able to connect to mongo : %s.", Failed, err)
			}
			t.Logf("\t%s\tShould be able to connect to mongo.", Success)
			defer dbConn.Close()

			cu := user.CreateUser{
				UserType:  1,
				FirstName: "bill",
				LastName:  "kennedy",
				Email:     "bill@ardanlabs.com",
				Company:   "ardan",
			}

			newUsr, err := user.Create(ctx, dbConn, &cu)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create user : %s.", Failed, err)
			}

			t.Logf("\t%s\tShould be able to create user.", Success)

			if _, err = user.Create(ctx, dbConn, &cu); err != nil {
				t.Fatalf("\t%s\tShould be able to create user : %s.", Failed, err)
			}

			t.Logf("\t%s\tShould be able to create user", Success)

			cu = user.CreateUser{
				UserType:  1,
				FirstName: "bill",
				LastName:  "smith",
				Email:     "bill@ardanlabs.com",
				Company:   "ardan",
			}

			if err := user.Update(ctx, dbConn, newUsr.UserID, &cu); err != nil {
				t.Fatalf("\t%s\tShould be able to update user : %s.", Failed, err)
			}
			t.Logf("\t%s\tShould be able to update user.", Success)

			rtv, err := user.Retrieve(ctx, dbConn, newUsr.UserID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user : %s.", Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve user.", Success)

			if rtv.LastName != cu.LastName {
				t.Errorf("\t%s\tShould be able to see updates to LastName.", Failed)
				t.Log("\t\tGot :", rtv.LastName)
				t.Log("\t\tWant :", cu.LastName)
			} else {
				t.Logf("\t%s\tShould be able to see updates to Lastname.", Success)
			}

			if err := user.Delete(ctx, dbConn, newUsr.UserID); err != nil {
				t.Fatalf("\t%s\tShould be able to delete user : %s.", Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete user.", Success)

			if _, err := user.Retrieve(ctx, dbConn, newUsr.UserID); err == nil {
				t.Fatalf("\t%s\tShould not be able to retrieve user : %s.", Failed, err)
			}
			t.Logf("\t%s\tShould not be able to retrieve user.", Success)

		}
	}
}

// Test set up
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

var ctx context.Context
var masterDB *db.DB

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	values := web.Values{
		TraceID: uuid.New(),
		Now:     time.Now(),
	}

	ctx = context.WithValue(context.Background(), web.KeyValues, &values)

	c, err := docker.StartMongo()
	if err != nil {
		log.Fatalln(err)
		docker.SetTestEnv(c)
	}

	docker.SetTestEnv(c)

	defer func() {
		if err := docker.StopMongo(c); err != nil {
			log.Println(err)
		}
	}()

	dbTimeout := 25 * time.Second
	dbHost := os.Getenv("DB_HOST")

	log.Println("main: started: Initialize Mongo")
	masterDB, err = db.New(dbHost, dbTimeout)
	if err != nil {
		log.Fatalf("startup : Register DB : %v", err)
	}

	defer masterDB.Close()

	return m.Run()
}
