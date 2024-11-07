# Student Management API

A RESTful API built with Go for managing student records with AI-powered summaries using Ollama.

## Features

- CRUD operations for student management
- Concurrent request handling
- In-memory data storage
- AI-powered student summaries using Ollama
- Input validation
- Error handling

## Prerequisites

- Go 1.x
- Ollama (for AI summaries)

## Installation

1. Clone the repository

```bash
git clone https://github.com/Pratik228/Students.git
cd students
```

2. Install dependencies

```bash
go mod tidy
```

3. Install Ollama and pull the required model

```bash
# Then pull the model:
ollama pull llama3.2
```

## Running the Application

1. Start the server:

```bash
go run main.go
```

2. The server will start on `http://localhost:8080`

## API Endpoints

- `POST /students` - Create a new student
- `GET /students` - Get all students
- `GET /students/{id}` - Get a student by ID
- `PUT /students/{id}` - Update a student
- `DELETE /students/{id}` - Delete a student
- `GET /students/{id}/summary` - Get AI-generated summary of a student

## Testing

Run the concurrency test:

```bash
go run tests/concurrency_load.go
```

## Postman Collection

Import the Postman collection using this link: [Postman Collection Link](https://lunar-robot-612579.postman.co/workspace/New-Team-Workspace~1f1a1d36-b6dd-4543-a51a-9c3a65603f9e/folder/19166721-dcd87862-d6bb-4bf5-a4f1-fe5fbe2a9ed9)

## API Usage Examples

### Create Student

```bash
POST http://localhost:8080/students
{
    "id": 101,
    "name": "John Doe",
    "age": 20,
    "email": "john@example.com"
}
```

### Get Student Summary

```bash
GET http://localhost:8080/students/101/summary
```
