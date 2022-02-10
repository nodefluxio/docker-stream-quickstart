-- migrate:up
ALTER TABLE "enrolled_face" 
  ADD COLUMN IF NOT EXISTS "gender" varchar(200),
  ADD COLUMN IF NOT EXISTS "birth_place" varchar(200),
  ADD COLUMN IF NOT EXISTS "birth_date" date;
-- migrate:down
ALTER TABLE "enrolled_face" 
  DROP COLUMN IF EXISTS "gender",
  DROP COLUMN IF EXISTS "birth_place",
  DROP COLUMN IF EXISTS "birth_date";
