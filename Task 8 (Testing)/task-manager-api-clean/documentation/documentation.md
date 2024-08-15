# Task Management API Documentation

## Endpoints

### Public Endpoints
- **POST /register** - Register a new user
- **POST /register-admin** - Register a new admin
- **POST /login** - Login user and get JWT token

### Protected Endpoints
- **POST /tasks** - Create a new task
- **GET /tasks** - Get all tasks
- **GET /tasks/:id** - Get a task by ID (only allowed to admin)
- **PUT /tasks/:id** - Update a task by ID (only allowed to admin)
- **DELETE /tasks/:id** - Delete a task by ID (only allowed to admin)
- **POST /promote/:id** - Promote a user to be an admin (only allowed to admin)

The endpoints for the `/tasks` API, as mentioned in the requirements, are implemented and tested using Postman. These are included in the published documentation.

## Authentication
Use JWT for authentication. Include the token in the `Authorization` header as `Bearer <token>` for protected routes.

## Features

### Database Integration
The API includes a MongoDB database integration for seamless data management.

### Authentication and Authorization
This version also includes implementations for authentication and authorization.

### Clean Architecture
The project follows clean architecture principles.

## Testing

### Automated Testing
The API uses automated tests to ensure the reliability and correctness of the endpoints. Here are the guidelines for automated testing:

#### Testing Framework
The API is tested using the Go testing package in conjunction with `testify/suite` to organize tests into suites.

#### Test Structure
Each endpoint and business logic is covered by tests that include:
- **Positive Tests:** Verifying correct responses to valid requests.
- **Negative Tests:** Checking error handling for invalid requests and edge cases.

#### Mocking Dependencies
Dependencies such as database connections and external services are mocked to isolate the code under test and focus on the business logic.

#### Running Tests
You can run the tests using the following command:
```bash
go test ./... 
```

### Manual Testing with Postman
For manual testing, Postman is used to validate API endpoints.

#### Postman Collection
The Postman collection includes pre-configured requests for all API endpoints. Each request is designed to test specific functionality.

#### Environment Variables
Use environment variables in Postman for dynamic values such as JWT tokens, user IDs, and task IDs.

#### Testing Scenarios
Ensure to test various scenarios including:
- Creating tasks with different data inputs.
- Fetching tasks with and without authentication.
- Updating and deleting tasks as an admin.
- Promoting users and testing role-based access.

#### Postman Documentation
You can find the updated Postman documentation at the following link:
[Postman Documentation](https://documenter.getpostman.com/view/33710823/2sA35HWLLr)

