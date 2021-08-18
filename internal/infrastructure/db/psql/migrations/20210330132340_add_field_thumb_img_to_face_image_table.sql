-- migrate:up
ALTER TABLE "face_image" ADD COLUMN IF NOT EXISTS "image_thumbnail" BYTEA;
-- migrate:down
ALTER TABLE "face_image" DROP COLUMN IF EXISTS "image_thumbnail";
