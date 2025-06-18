alter table "public"."ingredients" alter column "unit" drop not null;
alter table "public"."ingredients" add column "unit" varchar;
