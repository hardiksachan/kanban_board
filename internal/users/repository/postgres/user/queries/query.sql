-- name: GetUserById :one
SELECT *
FROM "user"
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM "user"
WHERE email = $1;

-- name: InsertUser :one
INSERT INTO "user"(name, email, password)
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateUser :one
UPDATE "user"
SET
    name = $1,
    email = $2,
    password = $3,
    modified_at = now()
WHERE
    user_id = $4
RETURNING *;

-- name: DeleteUser :one
DELETE FROM "user"
WHERE user_id = $1
RETURNING *;


