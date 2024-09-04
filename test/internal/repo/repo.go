/***
-- FILE
@file repository.go
@desc Defines the repository layer for user-related database operations. Provides methods to interact with the `users` table in the database.
@auth John Smith
@v 1.0
@date 01/01/2024
*/

package repository

/***
-- PKG
@pkg repository
@desc Provides the repository layer for user-related database operations. This package contains methods for interacting with the `users` table in the database, including retrieving user data.
@usage This package is used to perform database operations related to users, such as fetching all users or retrieving a specific user by ID. It is designed to interact with the database through the `UserRepository` type.
@dep (Model) Depends on the `model` package for user data representation.
*/

import (
	"database/sql"
	"myapp/internal/model"
)

/***
-- TYPE
@type UserRepository
@desc Repository for user-related database operations. Provides methods to retrieve user data from the `users` table.
@fields
	@field db (*sql.DB): Database connection used for executing SQL queries.
*/

type UserRepository struct {
	db *sql.DB
}

/***
-- FUNC
@func NewUserRepository
@desc Creates a new UserRepository instance with a given database connection.
@param dbConn (*sql.DB): Database connection to initialize the UserRepository.
@return (*UserRepository) Initialized UserRepository instance.
*/

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{
		db: dbConn,
	}
}

/***
-- FUNC
@func GetAllUsers
@desc Retrieves all users from the database.
@return ([]model.User) Slice of user models representing all users in the database.
@return (error) Any error encountered during the query execution.
@rec UserRepository
*/

func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

/***
-- FUNC
@func (r *UserRepository) GetUserByID
@desc Retrieves a user from the database by their ID.
@param id (int): ID of the user to retrieve.
@return (model.User) User model representing the user with the given ID.
@return (error) Any error encountered during the query execution or if the user is not found.
*/

func (r *UserRepository) GetUserByID(id int) (model.User, error) {
	var user model.User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	return user, err
}
