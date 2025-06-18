import GET_LIKE_STATUS from "../graphql/query/isLiked.gql";
import LIKE_RECIPE from "../graphql/mutation/recipes/likeRecipe.gql";
import REMOVE_LIKE from "../graphql/mutation/recipes/unlikeRecipe.gql";
import GET_RECIPE from "../graphql/query/recipe.gql";

export const useLikeStore = defineStore({
  id: "like",
  state: () => ({
    likes: [],
    isLiked: false,
    likeId: "",
    like: [],
    onload: false,
    bookmark: [],

    errorMessage: "",
    successmessage: "",

    processResultStatus: false,
  }),
  actions: {
    initRecipesState() {
      this.likes = [];
      this.onload = false;
      this.errorMessage = "";
      this.successmessage = "";
    },

    setErrorMessage(payload) {
      this.errorMessage = payload;
    },

    async likeRecipe(payload) {
      this.onload = true;
      const { recipe_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: LIKE_RECIPE,
          variables: {
            recipe_id,
          },
          awaitRefetchQueries: true,
          refetchQueries: [
            {
              query: GET_RECIPE,
              variables: {
                recipe_id,
              },
            },
          ],
          onCompleted: (data) => {
            console.log("the refetched recipe data", data);
          },
          onError: (error) => {
            console.log("error removing like", error);
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
        console.log(
          "like response from like store",
          JSON.stringify(res, null, 2)
        );
        this.isLiked = res.data.like.length > 0;
        this.likeId = res.data.like[0].id;
        console.log("id of the like", this.likeId);
        console.log("recipe liked");
      } catch (error) {
        console.log("error likeing the recipe");
      } finally {
        this.onload = false;
      }
    },

    async removeLike(payload) {
      this.onload = true;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: REMOVE_LIKE,
          variables: {
            id: payload,
          },
          awaitRefetchQueries: true,
          refetchQueries: [
            {
              query: GET_RECIPE,
              variables: {
                recipe_id,
              },
            },
          ],
          onCompleted: (data) => {
            console.log("the refetched recipe data", data);
          },
          onError: (error) => {
            console.log("error removing like", error);
          },
        });
        if (res) {
          this.successmessage = "recipe like removed";
          this.processResultStatus = true;
        }
      } catch (error) {
        this.processResultStatus = false;
        console.log("error removing like", error);
      }
    },
    setSuccessMessage(payload) {
      this.successmessage = payload;
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useLikeStore, import.meta.hot));
}
