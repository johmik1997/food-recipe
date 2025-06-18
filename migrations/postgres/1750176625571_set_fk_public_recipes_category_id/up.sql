alter table "public"."recipes" drop constraint "recipes_catagory_id_fkey",
  add constraint "recipes_category_id_fkey"
  foreign key ("category_id")
  references "public"."catagories"
  ("id") on update no action on delete no action;
