package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"

	"github.com/njdaniel/cassandra_example/cassandra"
)

// Get users
func Get(w http.ResponseWriter, r *http.Request) {
	var userList []User
	m := map[string]interface{}{}

	query := "SELECT id,age,firstname,lastname,city,email FROM users"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		userList = append(userList, User{
			ID:        m["id"].(gocql.UUID),
			Age:       m["age"].(int),
			FirstName: m["firstname"].(string),
			LastName:  m["lastname"].(string),
			Email:     m["email"].(string),
			City:      m["city"].(string),
		})
		m = map[string]interface{}{}
	}

	json.NewEncoder(w).Encode(AllUsersResponse{Users: userList})
}

// GetOne user
func GetOne(w http.ResponseWriter, r *http.Request) {
	var user User
	var errs []string
	var found = false

	vars := mux.Vars(r)
	id := vars["user_uuid"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id,age,firstname,lastname,city,email FROM users WHERE id=? LIMIT 1"
		iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			user = User{
				ID:        m["id"].(gocql.UUID),
				Age:       m["age"].(int),
				FirstName: m["firstname"].(string),
				LastName:  m["lastname"].(string),
				Email:     m["email"].(string),
				City:      m["city"].(string),
			}
		}
		if !found {
			errs = append(errs, "User not found")
		}
	}

	if found {
		json.NewEncoder(w).Encode(GetUserResponse{User: user})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}

// Enrich takes slice of UUIDs for db lookup and creates map of
// UUID and user's concat first and last name
// If empty list of UUIDs is passed, empty map is returned
func Enrich(uuids []gocql.UUID) map[string]string {
	if len(uuids) > 0 {
		names := map[string]string{}
		m := map[string]interface{}{}

		query := "SELECT id,firstname,lastname FROM users WHERE id IN ?"
		iterable := cassandra.Session.Query(query, uuids).Iter()
		for iterable.MapScan(m) {
			fmt.Println("m", m)
			userID := m["id"].(gocql.UUID)
			names[userID.String()] = fmt.Sprintf("%s %s", m["firstname"].(string), m["lastname"].(string))
			m = map[string]interface{}{}
		}
		return names
	}
	return map[string]string{}
}
