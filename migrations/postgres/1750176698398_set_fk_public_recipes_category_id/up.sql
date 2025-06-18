alter table "public"."recipes" drop constraint "recipes_category_id_fkey",
  add constraint "recipes_category_id_fkey"
  foreign key ("category_id")
  references "public"."categories"
  ("id") on update no action on delete no action;
