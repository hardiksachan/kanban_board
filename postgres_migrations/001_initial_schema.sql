CREATE TABLE IF NOT EXISTS "user"
(
    user_id     UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name        VARCHAR(50)      NOT NULL,
    email       VARCHAR(255)     NOT NULL,
    password    TEXT             NOT NULL,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT now(),
    modified_at TIMESTAMPTZ      NOT NULL DEFAULT now()
);

---- create above / drop below ----

DROP TABLE IF EXISTS "user" CASCADE;
