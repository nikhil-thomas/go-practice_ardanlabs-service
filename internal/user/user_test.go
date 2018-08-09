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

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/db"
)

func TestCreate(t *testing.T) {
	t.Log("Given the need to validate creating a user")
	{
		t.Log("\t using this user")
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

			usr, err := user.Create(ctx, dbConn, &cu)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create user : %s.", Failed, err)
			}

			t.Logf("\t%s\tShould be able to create user.", Success)

			const usersCollection = "users"
			q := bson.M{"user_id": usr.UserID}
			var u *user.User

			f := func(collection *mgo.Collection) error {
				return collection.Find(q).One(&u)
			}

			if err := dbConn.Execute(ctx, usersCollection, f); err != nil {
				if err == mgo.ErrNotFound {
					t.Fatalf("\t%s\tShould be able to see user in DB : %s.", Failed, err)
				}
				t.Logf("\t%s\tShould be able to see user in DB", Success)
			}
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
