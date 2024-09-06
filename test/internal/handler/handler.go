/***
-- FILE
@file handler.go
@desc Defines HTTP handlers for user-related endpoints, utilizing the service layer to process requests and interact with the database.
@auth John Smith
@v 1.0
@date 01/01/2024
*/

package handler

/***
-- PKG
@pkg handler
@desc Contains HTTP handlers for managing user-related endpoints. These handlers interact with the service layer to process requests and fetch or manipulate user data.
@usage This package is used to define routes and handlers for user operations such as retrieving user details and listing all users.
@dep (Testify) Used to test the handler functionality.
*/

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"myapp/internal/service"
	"net/http"
	"strconv"
)

/***
-- TYPE
@type UserHandler
@desc Handler for user-related HTTP requests, utilizing the user service to handle business logic.
@field service (UserService): Service for managing user-related operations.
@field service2 (Type2): This is here for testing.
@field service3 (Type3): This is here for testing.
*/

/***
-- VAR
@var ExampleVar
@type int
@desc This is a test var for this pkg.
*/

type UserHandler struct {
	service service.UserService
}

/***
-- FUNC
@func NewUserHandler
@desc Creates a new UserHandler instance with a given database connection.
@param dbConn (*sql.DB): Database connection to initialize the UserService.
@return (*UserHandler) Initialized UserHandler instance.
*/

func NewUserHandler(dbConn *sql.DB) *UserHandler {
	return &UserHandler{
		service: service.NewUserService(dbConn),
	}
}

/***
-- FUNC
@func GetAllUsers
@desc Handles HTTP GET requests to retrieve all users.
@rec UserHandler
*/

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

/***
-- FUNC
@func (h *UserHandler) GetUserByID
@desc Handles HTTP GET requests to retrieve a user by their ID.
@res (200) JSON encoded user object.
@res (400) If the provided user ID is invalid.
@res (404) If the user with the given ID does not exist.
*/

/***
-- VAR
@var exampleVar
@type int
@desc This is a test var for this pkg.
*/

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
