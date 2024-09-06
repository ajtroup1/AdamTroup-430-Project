/***
-- FILE
@file main.go
@desc Initializes the database connection, sets up the HTTP server, and routes requests to the handlers.
@auth John Smith
@v 1.2
@date 01/01/2024
*/

package main

/***
-- PKG
@pkg main
@desc Contains the high-level calls to <u>all</u> functionality in the app
@usage Entry point of the program. Simply 'run' the Makefile, and runtime starts here.
@dep (Gorilla Mux) Handles basic routing abilities necessary to register endpoints and serve them.
*/

import (
	"github.com/gorilla/mux"
	"log"
	"myapp/internal/handler"
	"myapp/pkg/db"
	"net/http"
)

/***
-- TYPE
@type testType
@desc This is a test for unexported type names.
@field field1 (Type): This is here for testing.
@field field2 (Type2): This is here for testing.
*/

/***
-- VAR
@var ExportedVar
@desc This is a test variable.
@type VariableType
*/

/***
-- FUNC
@func main
@desc The main function for the entire program. Creates a new handler using 'handler' and Gorilla Mux to listen and serve on port 8080. The example exists for testing purposes.
@param testParam (int): This is only here for testing.
@ret (string): This is only here for testing.
*/

func main() {
	// Initialize database connection
	conn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	/***
	-- VAR
	@var r
	@desc Gorilla Mux router. Via corresponding dependency
	@type *mux.Router
	*/

	// Initialize router
	r := mux.NewRouter()
	userHandler := handler.NewUserHandler(conn)
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")

	log.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", r)
}
