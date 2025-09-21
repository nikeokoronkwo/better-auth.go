package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

var authClient AuthClient

func init() {
	authClient = InitialiseClient(
		AuthClientOptions{
			Database: sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/postgres"),
		},
	)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", healthCheckHandler)

	authRouter := router.PathPrefix("/api/auth").Subrouter() 
	authRouter.HandleFunc("/sign-up/email", signUpHandler).Methods("POST")
	authRouter.HandleFunc("/sign-in/email", signInHandler).Methods("POST")
	authRouter.HandleFunc("/sign-out", signInHandler).Methods("POST")

	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}

// middleware function to ensure the given route
// is authorized
func authorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// run the func at the end
		handler(w, r)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Server Active: %s, ", r.Method))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	// read body
	var body SignUpEmailOptions
	bodyErr := json.NewDecoder(r.Body).Decode(&body);

	if bodyErr != nil {
		log.Fatalf("Failed to decode request body :: %v", bodyErr)
	}

	// kalel no
	result, err := authClient.SignUpEmail(body)
	if err != nil {
		log.Fatalf("Failed to sign up user :: %v", err)
	}

 	response, err := result.AsResponse()
	if err != nil {
		log.Fatalf("Failed to convert result to response :: %v", err)
	}

	response.Write(w)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	// read body
	var body SignInEmailOptions
	bodyErr := json.NewDecoder(r.Body).Decode(&body);

	if bodyErr != nil {
		log.Fatalf("Failed to decode request body :: %v", bodyErr)
	}

	// kalel no
	result, err := authClient.SignInEmail(body, r.Header)
	if err != nil {
		log.Fatalf("Failed to sign up user :: %v", err)
	}

 	response, err := result.AsResponse()
	if err != nil {
		log.Fatalf("Failed to convert result to response :: %v", err)
	}

	response.Write(w)
}

func signOutHandler(w http.ResponseWriter, r *http.Request) {
	// you need to give me the stone
	// (confused)
	result, err := authClient.SignOut(r.Header, struct{OnSuccess func()}{})
	if err != nil {
		log.Fatalf("Failed to sign up user :: %v", err)
	}

 	response, err := result.AsResponse()
	if err != nil {
		log.Fatalf("Failed to convert result to response :: %v", err)
	}

	response.Write(w)
}
