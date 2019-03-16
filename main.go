package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"redtoad.co/cassandra_example/cassandra"
	"redtoad.co/cassandra_example/users"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func main() {
	/*
		GET / — heartbeat check if our API is online

		GET /users — fetch all users from the database
		GET /users/UUID — fetch an individual user from the database
		POST /users/new — create a new user

		GET /messages — fetch all messages from Stream (with the database as a backup)
		GET /messages/UUID — fetch an individual message from the database
		POST /messages/new — create a new message
	*/
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)
	router.HandleFunc("/users/new", users.Post)
	router.HandleFunc("/users", Users.Get)
	router.HandleFunc("/users/{user_uuid}", Users.GetOne)
	log.Fatal(http.ListenAndServe(":8080", router))

}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
