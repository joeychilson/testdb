-- name: CreatePerson :exec
INSERT INTO people 
    (first_name, last_name, full_name, age, salary, start_date, phone, languages)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8);