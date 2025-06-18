import CREATE_BOOKMARK from "../graphql/mutation/bookmark/addBookmark.gql";
import GET_BOOKMARK_STATUS from "../graphql/query/isMarked.gql";
import REMOVE_BOOKMARK from "../graphql/mutation/bookmark/deleteBookmark.gql";
import GET_LIKE_STATUS from "../graphql/query/isLiked.gql";
import LIKE_RECIPE from "../graphql/mutation/recipes/likeRecipe.gql";
export const useBookmarkStore = defineStore({
  id: "bookmark",
  state: () => ({
    bookmarks: [],
    likes: [],
    isLiked: [],
    likeId: "",
    onload: false,
    bookmark: [],
    bookmarkedId: "",
    errorMessage: "",
    successmessage: "",
    isBookmarked: false,

    processResultStatus: false,
  }),
  actions: {
    initRecipesState() {
      this.bookmark = [];
      this.onload = false;
      this.bookmark = [];
      this.bookmarks = [];
      this.errorMessage = "";
      this.successmessage = "";
    },

    setErrorMessage(payload) {
      this.errorMessage = payload;
    },

    async createBookmark(payload) {
      this.onload = true;
      const { recipe_id } = payload;

      try {
        const { $apollo } = useNuxtApp();

        // Log payload before mutation
        console.log("Creating bookmark with payload:", payload);

        const res = await $apollo.clients.default.mutate({
          mutation: CREATE_BOOKMARK,
          variables: {
            recipe_id,
            
          },
        });

        // Log response to see if mutation was successful
        console.log("Bookmark response:", res);

        if (res.data && res.data.insert_bookmarks_one) {
          this.bookmarks.push(res.data.insert_bookmarks_one);
          this.processResultStatus = true;
        } else {
          throw new Error("Failed to save the bookmark.");
        }
      } catch (error) {
        console.log("Error saving a recipe:", error);
        this.errorMessage = error.message || "Error saving a recipe";

        this.processResultStatus = false;
      } finally {
        this.onload = false;
        return this.processResultStatus;
      }
    },

    async checkIfBookmarked(payload) {
      this.onload = true;
      const { recipe_id, user_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: GET_BOOKMARK_STATUS,
          variables: {
            recipe_id,
            user_id,
          },
        });

        this.isBookmarked = res.data.bookmarks.length > 0;
        console.log("is bookmarked", this.isBookmarked);
        this.bookmarkedId = res.data.bookmarks[0].id;
        console.log("bookmark id that is saved", this.bookmarkedId);
      } catch (error) {
        console.log("Error fetching bookmark status", error);
      } finally {
        this.onload = false;
      }
    },
    async removeBookmark(payload) {
      this.onload = true;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: REMOVE_BOOKMARK,
          variables: {
            id: payload,
          },
        });
        if (res) {
          this.successmessage = "Recipe unsaved!";
          this.processResultStatus = true;
        }
      } catch (error) {
        this.processResultStatus = false;
        console.log("Error removing bookmark:", error);
      }
    },

    async likeBookmark(payload) {
      this.onload = true;
      const { recipe_id, user_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: LIKE_RECIPE,
          variables: {
            recipe_id,
            user_id,
          },
        });
        console.log("like recipe response", res);
        if (res.data && res.data.insert_like_one) {
          this.likes.push(res.data.insert_like_one);
          this.processResultStatus = true;
        } else {
          console.log("error liking  arecipe");
        }
      } catch (error) {
        console.log("error liking recipe", error);
        this.errorMessage = error.message || "Error liking a recipe";
        this.processResultStatus = false;
      } finally {
        this.onload = false;
        return this.processResultStatus;
      }
    },

    async checkIfLiked(payload) {
      this.onload = true;
      const { recipe_id, user_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: GET_LIKE_STATUS,
          variables: {
            recipe_id,
            user_id,
          },
        });
        this.isLiked = res.data.like.length > 0;
        this.likeId = res.data.like.id;
        console.log("id of the like", this.likeId);
        console.log("recipe liked");
      } catch (error) {
        console.log("error likeing the recipe");
      } finally {
        this.onload = false;
      }
    },
    setSuccessMessage(payload) {
      this.successmessage = payload;
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useBookmarkStore, import.meta.hot));
}
