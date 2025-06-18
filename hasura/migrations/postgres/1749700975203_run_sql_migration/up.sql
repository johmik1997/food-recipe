-- Step 1: Alter column type and set NOT NULL
ALTER TABLE "public"."ratings"
ALTER COLUMN "rating" TYPE INT,
ALTER COLUMN "rating" SET NOT NULL;

-- Step 2: Add a CHECK constraint for allowed range
ALTER TABLE "public"."ratings"
ADD CONSTRAINT "rating_range_check"
CHECK ("rating" BETWEEN 1 AND 5);
