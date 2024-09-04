/***
-- FILE
@file service.go
@desc Defines the service layer for user-related operations. Provides methods to interact with the user repository and handle business logic.
@auth John Smith
@v 1.0
@date 01/01/2024
*/

package service

/***
-- PKG
@pkg service
@desc Contains the service layer for user-related operations. This package provides business logic and interacts with the `repository` package to manage user data. It offers methods to retrieve user information and perform operations related to users.
@usage This package is used to handle business logic for user operations, such as fetching all users or retrieving a specific user by ID. It communicates with the repository layer to access and manipulate user data.
@dep (Repository) Depends on the `repository` package for accessing user-related data from the database.
*/

import (
	"database/sql"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

/***
-- FUNC
@func NewUserService
@desc Creates a new UserService instance with a given database connection.
@param dbConn (*sql.DB):Database connection to initialize the UserRepository.
@return (*UserService) Initialized UserService instance.
*/

func NewUserService(dbConn *sql.DB) *UserService {
	return &UserService{
		repo: repository.NewUserRepository(dbConn),
	}
}

/***
-- FUNC
@func GetAllUsers
@desc Retrieves all users by calling the user repository.
@return ([]model.User) Slice of user models representing all users in the database.
@return (error) Any error encountered while retrieving users.
@rec UserService
*/

func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAllUsers()
}

/***
-- FUNC
@func (s *UserService) GetUserByID
@desc Retrieves a user by their ID by calling the user repository.
@param id (int): ID of the user to retrieve.
@return (model.User) User model representing the user with the given ID.
@return (error) Any error encountered while retrieving the user or if the user is not found.
*/

func (s *UserService) GetUserByID(id int) (model.User, error) {
	return s.repo.GetUserByID(id)
}
