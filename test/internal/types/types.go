/***
-- FILE
@file types.go
@desc Defines data types used throughout the application, including the user model with fields for user information. This description also contains the word package and pkg for testing reasons.
@auth John Smith
@v 1.0
@date 01/01/2024
*/

package types

/***
-- TYPE
@type User
@desc Represents a user in the application. This type includes fields for storing user ID, name, and email.
@fields
	@field ID (int): Unique identifier for the user.
	@field Name (string): Name of the user.
	@field Email (string): Email address of the user.
@usage This type is used to model user data for various operations including data storage and retrieval.
*/

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
