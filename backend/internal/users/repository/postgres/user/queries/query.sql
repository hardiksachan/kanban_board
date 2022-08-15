-- name: FindById :one
SELECT user_id, email, password
FROM "user"
WHERE user_id = $1;

-- name: FindByEmail :one
SELECT user_id, email, password
FROM "user"
WHERE email = $1;

-- name: CountByEmail :one
SELECT COUNT(*)
FROM "user"
WHERE email = $1;

-- name: InsertCredential :one
INSERT INTO "user"(email, password, name)
VALUES ($1, $2, split_part($1, '@', 1))
RETURNING user_id, email, password;

-- name: UpdatePassword :one
UPDATE "user"
SET password    = $1,
    modified_at = now()
WHERE user_id = $2
RETURNING user_id, email, password;

-- name: DeleteUser :one
DELETE
FROM "user"
WHERE user_id = $1
RETURNING *;

-- name: UpdateUserData :one
UPDATE "user"
SET name = $1,
    profile_image_url = $2,
    modified_at = now()
WHERE user_id = $3
RETURNING user_id, email, name, profile_image_url;

-- name: GetUserData :one
SELECT user_id, email, name, profile_image_url
FROM "user"
WHERE user_id = $1;