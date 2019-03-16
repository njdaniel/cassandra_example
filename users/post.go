package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/njdaniel/cassandra_example/cassandra"
)

// Post takes a gorilla/mux req and resp
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUUID gocql.UUID

	// FormToUser() is included in Users/processing.go
	// we will describe this later
	user, errs := FormToUser(r)

	// have we created a user correctly
	var created = false

	// if we had no errors from FormToUser, we will
	// attempt to save our data to Cassandra
	if len(errs) == 0 {
		fmt.Println("creating a new user")

		// generate a unique UUID for this user
		gocqlUUID = gocql.TimeUUID()

		// write data to Cassandra
		if err := cassandra.Session.Query(`
		INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
			gocqlUUID, user.FirstName, user.LastName, user.Email, user.City, user.Age).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	// depending on whether we created the user, return the
	// resource ID in a JSON payload, or return our errors
	if created {
		fmt.Println("user_id", gocqlUUID)
		json.NewEncoder(w).Encode(NewUserResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
