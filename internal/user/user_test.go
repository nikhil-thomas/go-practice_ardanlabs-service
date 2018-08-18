package user_test

import (
	"os"
	"testing"

	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/platform/tests"
	"github.com/nikhil-thomas/go-practice_ardanlabs-service/internal/user"
)

var test *tests.Test

// TestMain is the entry point for testing
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	test = tests.New()
	defer test.TearDown()
	return m.Run()
}

func TestUser(t *testing.T) {
	defer tests.Recover(t)
	t.Log("Given the need to validate creating a user")
	{
		t.Log("\t when handling a single user")
		{
			ctx := tests.Context()
			dbConn := test.MasterDB.Copy()
			t.Logf("\t%s\tShould be able to connect to mongo.", tests.Success)
			defer dbConn.Close()

			cu := user.CreateUser{
				UserType:  1,
				FirstName: "bill",
				LastName:  "kennedy",
				Email:     "bill@ardanlabs.com",
				Company:   "ardan",
				Addresses: []user.CreateAddress{
					{
						Type:    1,
						LineOne: "12973 SW 112th ST",
						LineTwo: "Suite 153",
						City:    "Miami",
						State:   "FL",
						Zipcode: "33172",
						Phone:   "305-527-3353",
					},
				},
			}

			newUsr, err := user.Create(ctx, dbConn, &cu)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create user : %s.", tests.Failed, err)
			}

			t.Logf("\t%s\tShould be able to create user.", tests.Success)

			usr, err := user.Retrieve(ctx, dbConn, newUsr.UserID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user : %s.", tests.Failed, err)
			}

			t.Logf("\t%s\tShould be able to retrieve user", tests.Success)

			cu = user.CreateUser{
				UserType:  1,
				FirstName: "bill",
				LastName:  "smith",
				Email:     "bill@ardanlabs.com",
				Company:   "ardan",
			}

			if err := user.Update(ctx, dbConn, newUsr.UserID, &cu); err != nil {
				t.Fatalf("\t%s\tShould be able to update user : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update user.", tests.Success)

			usr, err = user.Retrieve(ctx, dbConn, newUsr.UserID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve user : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve user.", tests.Success)

			if usr.LastName != cu.LastName {
				t.Log("\t\tGot :", usr.LastName)
				t.Log("\t\tWant :", cu.LastName)
				t.Errorf("\t%s\tShould be able to see updates to LastName.", tests.Failed)
			} else {
				t.Logf("\t%s\tShould be able to see updates to Lastname.", tests.Success)
			}

			if err := user.Delete(ctx, dbConn, usr.UserID); err != nil {
				t.Fatalf("\t%s\tShould be able to delete user : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete user.", tests.Success)

			usr, err = user.Retrieve(ctx, dbConn, newUsr.UserID)
			if err != user.ErrNotFound {
				t.Fatalf("\t%s\tShould not be able to retrieve user : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould not be able to retrieve user.", tests.Success)
		}
	}
}
