# Chirpy

A small web server application built with go 1.23 (https://go.dev/doc/install) and PostgreSQL 17.4. This is a guided project, following steps from course [Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang) provided by [Boot.dev](https://www.boot.dev/). This backend application allows users to create and delete Chirps (text posts with a limit of 140 characters).

## Goal

- Understand how to create a web server using Go's standard library
- Better understand RESTful API design and best practices
- Implement secure authentication/authorization
- Learn about webhooks

## How to run

Ensure you have Go 1.23 and Postgres 17.4 installed. Once complete, clone this repository and run the command `go mod tidy` from the root directory in this project. This project relies on the dependency [godotenv](https://github.com/joho/godotenv). Add a `.env` file to the root directory of this project and give it these fields:

```
DB_URL="<Your database URL>?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="<Your JSON Web Token secret>"
POLKA_KEY="<An api key used in authorization header of calls made to the webhooks endpoint>"
```

To generate your own JWT_SECRET you can run the command `openssl rand -base64 64` from your terminal. Similarly, you can generate a POLKA_KEY with the same command, just with 32 characters (i.e. `openssl rand -base64 32`). From there, open up a new terminal from the root directory and run either `go run .` or `go build -o out && ./out`. The latter command will generate the binary file in the root directory and run it. If the application started successfully, you will be able to see it by opening a browser and navigating to `localhost:8080/app/`. You can also navigate to `localhost:8080/admin/metrics` to view how many times the homepage has ben hit.

## API Documentation

There are several API endpoints accessible to this web server. Below are the requirements and expected outputs of each endpoint

### Health Check

#### GET /api/healthz

Check the status of the API

Response:
`Status 200 OK`

### Metrics

#### GET /admin/metrics

Response:

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 1 times!</p>
  </body>
</html>
```

### Users

#### POST /api/users

Register new user

Request body:

```json
{
  "email": "example@test.com",
  "password": "Ex4mple!"
}
```

Response:
`Status: 201 Created`

```json
{
  "id": "fd8f3194-5af4-47ce-bbf3-d810351512dd",
  "created_at": "2025-04-09T15:27:56.20467Z",
  "updated_at": "2025-04-09T15:27:56.20467Z",
  "email": "example@test.com",
  "is_chirpy_red": false
}
```

#### PUT /api/users

Update existing user

Header required:
`Authorization: Bearer <JWT>`

Request body required:

```json
{
  "email": "example@test.com",
  "password": "Ex4mple!"
}
```

Response:
`Status: 200 OK`

```json
{
  "id": "fd8f3194-5af4-47ce-bbf3-d810351512dd",
  "created_at": "2025-04-09T15:27:56.20467Z",
  "updated_at": "2025-04-09T15:35:34.436396Z",
  "email": "example@test.com",
  "is_chirpy_red": false
}
```

### Login

#### POST /api/login

Login with email and password

Request body required:

```json
{
  "email": "example@test.com",
  "password": "Ex4mple!"
}
```

Response:
`Status: 200 OK`

```json
{
  "id": "fd8f3194-5af4-47ce-bbf3-d810351512dd",
  "created_at": "2025-04-09T15:27:56.20467Z",
  "updated_at": "2025-04-09T15:35:34.436396Z",
  "email": "example@test.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJmZDhmMzE5NC01YWY0LTQ3Y2UtYmJmMy1kODEwMzUxNTEyZGQiLCJleHAiOjE3NDQyMzgzMzIsImlhdCI6MTc0NDIzNDczMn0.AnuT-gMqZGtqmkWuDa6FF0QMQHtfwpHDGnm6y-3eCtc",
  "refresh_token": "2f8a976ad34540bcd89c49edd124f58b805a93bc13f096bd3dd1036609a99ae5"
}
```

### Chirps

#### POST /api/chirps

Create a new Chirp

Header required:
`Authorization: Bearer <JWT>`

Request body required:

```json
{
  "body": "Chirp message",
  "user_id": "02320105-abd3-4ec7-adea-57e5d838d21c"
}
```

Response:
`Status: 201 Created`

```json
{
  "id": "a797bb2e-eb54-4855-93e9-2b0cebfb3986",
  "created_at": "2025-04-09T15:56:40.092149Z",
  "updated_at": "2025-04-09T15:56:40.092149Z",
  "body": "Chirp message",
  "user_id": "fd8f3194-5af4-47ce-bbf3-d810351512dd"
}
```
