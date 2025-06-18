import CREATE_RATING from "../graphql/mutation/rating/insertRating.gql";
import GET_RECIPE from "../graphql/query/recipe.gql"; //
export const useRatingStore = defineStore({
  id: "rating",
  state: () => ({
    ratings: [],
    onload: false,
    rating: [],
    errorMessage: "",
    successmessage: "",

    processResultStatus: false,
  }),
  actions: {
    initRecipesState() {
      this.ratings = [];
      this.onload = false;
      this.rating = [];
      this.errorMessage = "";
      this.successmessage = "";

      this.processResultStatus = false;
    },

    setErrorMessage(payload) {
      this.errorMessage = payload;
    },

    async createRating(payload) {
      this.onload = true;
      const { rating, recipe_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: CREATE_RATING,
          variables: {
            rating,
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
          onCompeleted: (data) => {
            console.log("the refetched recipe data", data);
          },
          onError: (error) => {
            console.log("errro refteching the recipe", error);
          },
        });
        this.rating = res.data.insert_ratings_one;
        this.processResultStatus = true;
      } catch (error) {
        console.log("error raing a recipe", error);
        this.errorMessage = error.message || "error rating a recipe";
        this.processResultStatus = false;
      } finally {
        this.onload = false;
        return this.processResultStatus;
      }
    },

    setSuccessMessage(payload) {
      this.successmessage = payload;
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useRatingStore, import.meta.hot));
}
