# GoDoc generator documentation
## Packages:
  - ### Package: `main`
    Contains the high-level calls to <u>all</u> functionality in the app

---
  - ### Package: `db`
    Contains functions for interacting with the database, specifically for establishing and managing connections.

      #### Functions for `db`:
      - **NewConnection**
        - Creates a new connection to the MySQL database using the provided Data Source Name (DSN).
        - 
        - Parameters:
---
  - ### Package: `handler`
    Contains HTTP handlers for managing user-related endpoints. These handlers interact with the service layer to process requests and fetch or manipulate user data.

      #### Types:
      - **UserHandler**
        - Handler for user-related HTTP requests, utilizing the user service to handle business logic.
        - Fields:
          - `service`
            - Data type: `UserService`
            - Service for managing user-related operations.
          - `service2`
            - Data type: `Type2`
            - This is here for testing.
          - `service3`
            - Data type: `Type3`
            - This is here for testing.
      #### Functions for `handler`:
      - **NewUserHandler**
        - Creates a new UserHandler instance with a given database connection.
        - 
        - Parameters:
          - `dbConn`
            - Data type: `*sql.DB`
            - Database connection to initialize the UserService.
      - **GetAllUsers**
        - Handles HTTP GET requests to retrieve all users.
        - UserHandler
        - Parameters:
---
  - ### Package: `repository`
    Provides the repository layer for user-related database operations. This package contains methods for interacting with the `users` table in the database, including retrieving user data.

      #### Types:
      - **UserRepository**
        - Repository for user-related database operations. Provides methods to retrieve user data from the `users` table.
        - Fields:
          - `db`
            - Data type: `*sql.DB`
            - Database connection used for executing SQL queries.
      #### Functions for `repository`:
      - **NewUserRepository**
        - Creates a new UserRepository instance with a given database connection.
        - 
        - Parameters:
          - `dbConn`
            - Data type: `*sql.DB`
            - Database connection to initialize the UserRepository.
      - **GetAllUsers**
        - Retrieves all users from the database.
        - UserRepository
        - Parameters:
---
  - ### Package: `service`
    Contains the service layer for user-related operations. This package provides business logic and interacts with the `repository` package to manage user data. It offers methods to retrieve user information and perform operations related to users.

      #### Functions for `service`:
      - **NewUserService**
        - Creates a new UserService instance with a given database connection.
        - 
        - Parameters:
          - `dbConn`
            - Data type: `*sql.DB`
            - Database connection to initialize the UserRepository.
      - **GetAllUsers**
        - Retrieves all users by calling the user repository.
        - UserService
        - Parameters:
---
  - ### Package: `types`
    Contains the types necessary for the entire program

      #### Types:
      - **User**
        - Represents a user in the application. This type includes fields for storing user ID, name, and email.
        - Fields:
          - `ID`
            - Data type: `int`
            - Unique identifier for the user.
          - `Name`
            - Data type: `string`
            - Name of the user.
          - `Email`
            - Data type: `string`
            - Email address of the user.
---
