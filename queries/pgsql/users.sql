-- name: CreateUser :one
INSERT INTO "users" ("first_name", "last_name", "email")
VALUES ($1, $2, $3)
RETURNING "id";