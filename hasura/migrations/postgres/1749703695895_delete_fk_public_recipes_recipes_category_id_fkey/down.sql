alter table "public"."recipes"
  add constraint "recipes_category_id_fkey"
  foreign key ("category_id")
  references "public"."categories"
  ("id") on update no action on delete set null;
