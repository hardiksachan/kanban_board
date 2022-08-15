ALTER TABLE "user"
    DROP COLUMN display_name;

---- create above / drop below ----


ALTER TABLE "user"
    ADD COLUMN display_name VARCHAR(50);

UPDATE "user"
SET display_name = name;

ALTER TABLE "user"
    ALTER COLUMN display_name SET NOT NULL;
