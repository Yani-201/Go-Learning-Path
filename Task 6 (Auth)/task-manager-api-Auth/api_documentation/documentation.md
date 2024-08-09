# Task Management API Documentation


## Endpoints

### Public Endpoints

- `POST /register` - Register a new user
- `POST /register-admin` - Register a new admin
- `POST /login` - Login user and get JWT token

### Protected Endpoints

- `POST /tasks` - Create a new task
- `GET /tasks` - Get all tasks
- `GET /tasks/:id` - Get a task by ID (only allowed to admin)
- `PUT /tasks/:id` - Update a task by ID (only allowed to admin)
- `DELETE /tasks/:id` - Delete a task by ID (only allowed to admin )
- `POST /promote/:id` - Promote a user to be an admin (only allowed to admin)


The endpoints for the `/tasks` API, as mentioned in the requirements, are implemented and tested using Postman. These are included in the published documentation.

## Authentication

- Use JWT for authentication. Include the token in the `Authorization` header as `Bearer <token>` for protected routes.


The documentation for the Task Management API is created and published using Postman.


## Features

- **Database Integration**: The API now includes a MongoDB database integration for seamless data management.
- **Authentication and Authorization**: This version also includes implementations for authentication and authorization.

## Postman Documentation

You can find the updated Postman documentation at the following link:  
[Postman Documentation](https://documenter.getpostman.com/view/33710823/2sA35HWLLr)