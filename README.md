# Todo App API

This is a simple API built with Go, Gin, and MongoDB to manage tasks in a todo list.

## Features
- Create tasks
- Read all tasks
- Update tasks
- Delete tasks

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/todoApp.git
    ```
2. Install dependencies:
    ```bash
    go mod tidy
    ```
3. Make sure MongoDB is running on `localhost:27017`.
4. Run the application:
    ```bash
    go run main.go
    ```

The API will be running on `http://localhost:8080`.

## API Endpoints

- `GET /tasks`: Get all tasks
- `POST /tasks`: Create a new task
- `PUT /tasks/{id}`: Update a task by ID
- `DELETE /tasks/{id}`: Delete a task by ID
