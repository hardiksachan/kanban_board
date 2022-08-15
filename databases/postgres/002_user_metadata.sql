ALTER TABLE "user"
    ADD COLUMN profile_image_url TEXT,
    ADD COLUMN display_name      VARCHAR(50);

UPDATE "user"
SET display_name = name;

ALTER TABLE "user"
    ALTER COLUMN display_name SET NOT NULL;

---- create above / drop below ----

ALTER TABLE "user"
    DROP COLUMN profile_image_url CASCADE;

ALTER TABLE "user"
    DROP COLUMN display_name CASCADE;
