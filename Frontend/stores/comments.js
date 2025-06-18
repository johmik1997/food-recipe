import CREATE_COMMENT from "../graphql/mutation/comment/insertComent.gql";
import UPDATE_COMMENT from "../graphql/mutation/comment/updateComment.gql";
import DELETE_COMMENT from "../graphql/mutation/comment/deleteComment.gql";
import GET_RECIPE from "../graphql/query/recipe.gql";
export const useCommnentStore = defineStore({
  id: "comments",
  state: () => ({
    comments: [],
    onload: false,
    comment: [],
    errorMessage: "",
    successmessage: "",

    processResultStatus: false,
  }),
  actions: {
    initRecipesState() {
      this.comments = [];
      this.onload = false;
      this.comment = [];
      this.errorMessage = "";
      this.successmessage = "";

      this.processResultStatus = false;
    },

    setErrorMessage(payload) {
      this.errorMessage = payload;
    },

    async createComment(payload) {
      this.onload = true;
      const { comment, recipe_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: CREATE_COMMENT,
          variables: {
            comment,
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
            console.log("error refetching the recipe details", error);
          },
        });
        this.comment = res.data.insert_comments_one;
        this.processResultStatus = true;
      } catch (error) {
        console.log("error creating comment", error);
        this.errorMessage = error.message || "error creating comment";
        this.processResultStatus = false;
      } finally {
        this.onload = false;
        return this.processResultStatus;
      }
    },

    async deleteComment(payload) {
      this.onload = true;
      const { id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: DELETE_COMMENT,
          variables: {
            id,
          },
        });
        this.comment = res.data.delete_comments;
        this.processResultStatus = true;
      } catch (error) {
        console.log("error deleting comment", error);
        this.errorMessage = error.message || "error deleting comment";
        this.processResultStatus = false;
      } finally {
        this.onload = false;
        return this.processResultStatus;
      }
    },

    async updateComment(payload) {
      this.onload = true;
      const { id, comment } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: UPDATE_COMMENT,
          variables: {
            id,
            comment,
          },
        });
        this.comment = res.data.update_comments;
        this.processResultStatus = true;
      } catch (error) {
        console.log("error updating comment", error);
        this.errorMessage = error.message || "error updating comment";
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
  import.meta.hot.accept(acceptHMRUpdate(useCommnentStore, import.meta.hot));
}
