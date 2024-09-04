/***
-- FILE
@file db.go
@desc Provides functions for establishing a database connection using the MySQL driver.
@auth John Smith
@v 1.0
@date 01/01/2024
*/

package db

/***
-- PKG
@pkg db
@desc Contains functions for interacting with the database, specifically for establishing and managing connections.
@usage This package provides the `NewConnection` function to create and return a new database connection.
@dep (MySQL Driver) Depends on the `github.com/go-sql-driver/mysql` package to interact with MySQL databases.
*/

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/***
-- FUNC
@func NewConnection
@desc Creates a new connection to the MySQL database using the provided Data Source Name (DSN).
@return (*sql.DB) Database connection instance.
@return (error) Any error encountered while opening the database connection.
*/

func NewConnection() (*sql.DB, error) {
	dsn := "user:password@tcp(127.0.0.1:3306)/mydatabase"
	return sql.Open("mysql", dsn)
}