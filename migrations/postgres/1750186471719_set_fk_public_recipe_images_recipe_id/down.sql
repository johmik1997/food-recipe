alter table "public"."recipe_images" drop constraint "recipe_images_recipe_id_fkey",
  add constraint "recipe_images_recipe_id_fkey"
  foreign key ("recipe_id")
  references "public"."recipes"
  ("id") on update cascade on delete cascade;
