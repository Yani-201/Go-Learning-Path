# Task Manager API Documentation

## Endpoints

### GET /tasks
- **Description**: Retrieves all tasks.
- **Response**: A list of all tasks.

### GET /tasks/:id
- **Description**: Retrieves a task by ID.
- **Response**: The task with the specified ID.

### POST /tasks
- **Description**: Creates a new task.
- **Request Body**: Task object.
- **Response**: Confirmation message.

### PATCH /tasks/:id
- **Description**: Updates a task by ID.
- **Request Body**: Fields to update.
- **Response**: Confirmation message.

### DELETE /tasks/:id
- **Description**: Deletes a task by ID.
- **Response**: Confirmation message.



The documentation for the Task Management API is created and published using Postman.

The endpoints for the `/tasks` API, as mentioned in the requirements, are implemented and tested using Postman. This is also included in the published documentation.

The API now includes a MongoDB database integration for seamless data management.

You can find the Postman Documentation at the following link: [Postman Documentation]
(https://documenter.getpostman.com/view/33710823/2sA35A74yU)