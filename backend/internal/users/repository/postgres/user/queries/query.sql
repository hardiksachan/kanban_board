-- name: GetUserById :one
SELECT user_id, name, email, password, created_at, modified_at
FROM "user"
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT user_id, name, email, password, created_at, modified_at
FROM "user"
WHERE email = $1;

-- name: InsertUser :one
INSERT INTO "user"(name, email, password, display_name)
VALUES ($1, $2, $3, $1)
RETURNING user_id, name, email, password,created_at, modified_at;

-- name: UpdateUser :one
UPDATE "user"
SET name        = $1,
    email       = $2,
    password    = $3,
    modified_at = now()
WHERE user_id = $4
RETURNING user_id, name, email, password, created_at, modified_at;

-- name: DeleteUser :one
DELETE
FROM "user"
WHERE user_id = $1
RETURNING user_id, name, email, password,created_at, modified_at;

-- name: UpdateUserMetadata :one
UPDATE "user"
SET display_name      = $1,
    profile_image_url = $2,
    modified_at       = now()
WHERE user_id = $3
RETURNING user_id, display_name, profile_image_url;

-- name: GetUserMetadata :one
SELECT user_id, display_name, profile_image_url
FROM "user"
WHERE user_id = $1;


