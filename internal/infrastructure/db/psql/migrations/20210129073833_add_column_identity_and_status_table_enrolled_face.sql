-- migrate:up
ALTER TABLE "enrolled_face" 
  ADD COLUMN IF NOT EXISTS "identity_number" varchar(200),
  ADD COLUMN IF NOT EXISTS "status" varchar(200);
-- migrate:down
ALTER TABLE "enrolled_face" 
  DROP COLUMN IF EXISTS "identity_number",
  DROP COLUMN IF EXISTS "status";
