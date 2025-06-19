import { acceptHMRUpdate, defineStore } from "pinia";
import RECIPES_QUERY from "../graphql/query/recipes.gql";
import GET_RECIPES from "../graphql/query/recipes/recipes.gql";
import GET_RECIPE from "../graphql/query/recipe.gql";
import UPLOAD_RECIPE_IMAGES from "../graphql/mutation/uploadRecipeImage.gql";
import CREATE_RECIPE from "../graphql/mutation/recipes/addRecipe.gql";
import CATEGORIES from "../graphql/query/catagories.gql";
import FILTER_BY_CATEGORY from "../graphql/query/filterByCategory.gql";
import UNLIKE_RECIPE from "../graphql/mutation/recipes/unlikeByRecipeid.gql";
import UNSAVE_RECIPE from "../graphql/mutation/recipes/unsaveRecipebyId.gql";
import BUY_RECIPE from "../graphql/mutation/recipes/buyRecipe.gql";
import VERIFY_PAYMENT from "../graphql/mutation/recipes/verifyPayment.gql";
import CHECKOUT_URL from "../graphql/query/checkoutUrl.gql";
import SELLER_RECIPES from "../graphql/query/soldRecipebySellerid.gql";
import DELETE_RECIPE from "../graphql/mutation/recipes/deleteRecipe.gql";
import PAYMENT_DETAILS from "../graphql/query/paymentDetail.gql";
import UPDATE_RECIPE from "../graphql/mutation/recipes/updateRecipe.gql";
import PROFILE_QUERY from "../graphql/query/profile.gql";

export const useRecipeStore = defineStore({
  id: "recipe",
  state: () => ({
    recipes: [],
    soldRecipes: [],
    categorizedRecipes: [],
    isLoading: false,
    recipe: [],
    errorMessage: "",
    successmessage: "",
    limit: 100,
    offset: 0,
    currentPage: 1,
    totalRecipes: 0,
    categories: [],
    createdRecipe: [],
    updatedRecipe: [],
    checkoutUrl: "",
    paymentId: "",
    isPaymentVerified: false,
    searchRecipe: "%%",
    processResultStatus: false,
    paymentStatus: "",
    sellerId: "",
    paymentDetails: [],
    buyerCount: 0,
    tx_ref: "",
  }),
  actions: {
    setOffset(payload) {
      this.offset = payload;
    },
    setLimit(payload) {
      this.limit = payload;
    },
    setCurrentPage(payload) {
      this.currentPage = Math.max(1, payload);
      this.offset = (this.currentPage - 1) * this.limit;
    },
    setErrorMessage(payload) {
      this.errorMessage = payload;
    },
    setSuccessMessage(payload) {
      this.successmessage = payload;
    },
    setSearchRecipe(payload) {
      this.searchRecipe = "%" + payload + "%";
    },

    async getAllRecipes(payload = "") {
      this.isLoading = true;

      const searchQuery = "%" + payload + "%";
      const { $apollo } = useNuxtApp();

      try {
        const res = await $apollo.clients.default.query({
          query: GET_RECIPES,
          variables: {
            limit: this.limit,
            offset: this.offset,
            searchQuery,
          },
          fetchPolicy: "network-only",
        });

        console.log("variables for the query from the store", {
          limit: this.limit,
          offset: this.offset,
          searchQuery: "%%",
        });

        this.recipes = res.data.recipes || [];
        this.totalRecipes = res.data.recipes_aggregate.aggregate.count;
        console.log("search query from the store", this.searchRecipe);
        console.log(
          "Fetched recipes from the store:",
          JSON.stringify(res.data.recipes, null, 2)
        );
      } catch (error) {
        console.error("Error loading recipes:", error);
        this.setErrorMessage("Failed to load recipes");
      } finally {
        this.isLoading = false;
      }
    },
    async resetSearch() {
      this.setSearchRecipe("");
      await this.getAllRecipes();
    },

    async singleRecipe(payload) {
      this.isLoading = true;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: GET_RECIPE,
          variables: {
            recipeId: payload,
          },
          fetchPolicy: "network-only",
        });
        console.log("recipeissss", res);
        this.recipe = res.data.recipes_by_pk;
        console.log("the fetched recipe data", this.recipe);
      } catch (error) {
        console.log("error fetching recipes", error);
      } finally {
        this.isLoading = false;
      }
    },

    async createRecipe(payload) {
      this.isLoading = true;
      const {
        title,
        description,
        prep_time,
        cook_time,
        servings,
        category_id,
        price,
        ingredients,
        user_id,
        steps,
        imageName,
        imageType,
        base64string,
      } = payload;

      console.log("paylooooood", JSON.stringify(payload));
      const searchQuery = "%" + payload + "%";
      try {
        const { $apollo } = useNuxtApp();
        let variables = {
          title,
          description,
          prep_time,
          cook_time,
          servings,
          category_id,
          price,
          ingredients,
          user_id,
          steps,
          base64String: base64string,
          type: imageType,
          name: imageName,
        };

        // Execute the mutation with await
        const res = await $apollo.clients.default.mutate({
          mutation: CREATE_RECIPE,
          variables,
          awaitRefetchQueries: true,
          refetchQueries: [
            {
              query: GET_RECIPES,
              variables: {
                limit: this.limit,
                offset: this.offset,
                searchQuery: "%%",
              },
            },
          ],
          onCompleted: (data) => {
            console.log("the refetched recipe data", data);
          },
          onError: (error) => {
            console.log("error refecthing the queries", error);
          },
        });

        // Log the full response for debugging
        console.log("Mutation response:", res);

        // Handle the response
        if (res.data && res.data.addRecipe) {
          const data = res.data.addRecipe;
          console.log("resultt", JSON.stringify(res.data.addRecipe));
          this.createdRecipe = data;
          console.log("created recipe", JSON.stringify(this.createdRecipe));
          this.successMessage = data.message || "Recipe created successfully";
          this.processResultStatus = true;
        } else {
          throw new Error("Invalid response from server");
        }
      } catch (error) {
        console.log("Error creating a recipeEEEE", error);
        this.errorMessage = error.message || "Couldn't create a recipe";
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
        return this.processResultStatus;
      }
    },

    async uploadRecipeImages(payload) {
      this.isLoading = true;
      const { recipe_id, images, featuredImageIndex } = payload;
      const { $apollo } = useNuxtApp();

      console.log("recipe id", recipe_id);
      console.log("images payload", images);

      try {
        const res = await $apollo.clients.default.mutate({
          mutation: UPLOAD_RECIPE_IMAGES,
          variables: {
            recipe_id,
            images,
            featuredImageIndex,
          },
        });

        console.log("Mutation response:", res);

        if (res.data && res.data.uploadRecipeImage) {
          this.successMessage = "Images uploaded successfully";
          this.processResultStatus = true;
        } else {
          throw new Error(res.errors?.[0]?.message || "Error uploading images");
        }
      } catch (error) {
        console.log("Error uploading image", error);
        this.errorMessage = error.message || "Error uploading images";
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
        return this.processResultStatus;
      }
    },

    async getCategories() {
      this.isLoading = true;
      const { $apollo } = useNuxtApp();
      try {
        this.processResultStatus = true;
        const res = await $apollo.clients.default.query({
          query: CATEGORIES,
        });
        if (res.data) {
          this.recipes = res.data.recipes;
          console.log("categories", this.recipe);
        }
        this.categories = res.data.categories;
      } catch (error) {
        this.processResultStatus = false;
        console.log("error fetching catagories", error);
      } finally {
        this.isLoading = false;
        return this.processResultStatus;
      }
    },

    async filterByCategory(payload = "") {
      const name = "%" + payload + "%";
      this.isLoading = true;

      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: FILTER_BY_CATEGORY,
          variables: {
            name,
            limit: this.limit,
            offset: this.offset,
          },
        });
        this.recipes = res.data.recipes || [];
      } catch (error) {
        console.error("Error loading recipes:", error);
      } finally {
        this.isLoading = false;
      }
    },

    async unlikeRecipe(payload) {
      this.isLoading = true;
      try {
        const { $apollo } = useNuxtApp();
        const { recipe_id, id } = payload;
        const res = await $apollo.clients.default.mutate({
          mutation: UNLIKE_RECIPE,
          variables: {
            recipe_id,
            id,
          },
          refetchQueries: [
            {
              query: PROFILE_QUERY,
              variables: {
                id,
              },
            },
          ],
          onCompleted: (data) => {
            console.log("the refetched recipe data", data);
          },
          onError: (error) => {
            console.log("error unliking a recipe", error);
          },
        });
        console.log("recipe id to unlike from store", payload);
        if (res) {
          this.successmessage = "recipe unliked successfully";
          this.processResultStatus = true;
        }
      } catch (error) {
        console.log("error unliking a recipeeeee", error);
        this.processResultStatus = false;
      } finally {
        return this.processResultStatus;
      }
    },

    async unsaveRecipe(payload) {
      this.isLoading = true;
      try {
        const { $apollo } = useNuxtApp();
        const { recipe_id, id } = payload;
        const res = await $apollo.clients.default.mutate({
          mutation: UNSAVE_RECIPE,
          variables: {
            recipe_id,
            id,
          },
          refetchQueries: [
            {
              query: GET_RECIPES,
              variables: {
                limit: this.limit,
                offset: this.offset,
                searchQuery: "%%",
              },
            },
          ],
        });
        console.log("recipe id to unsave delete from store", payload);
        if (res) {
          this.successmessage = "recipe unsaved successfully";
          this.processResultStatus = true;
        }
      } catch (error) {
        console.log("error unsaving a recipeeeee", error);
        this.processResultStatus = false;
      } finally {
        return this.processResultStatus;
      }
    },

    async buyRecipe(payload) {
      this.isLoading = true;
      const { buyer_id, recipe_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.mutate({
          mutation: BUY_RECIPE,
          variables: {
            buyer_id,
            recipe_id,
          },
        });
        if (res.data) {
          const data = res.data.buyRecipe;
          console.log(
            "result from buying reciepe store",
            JSON.stringify(data, null, 2)
          );
          this.paymentId = data.payment_id;
          console.log("payment id from store", this.paymentId);
          this.checkOutUrl = data.checkout_url;
          console.log("checkout url from store", this.checkOutUrl);
          this.sellerId = data.seller_id;
          console.log("seller id from store", this.sellerId);
          this.successmessage = data.message || "recipe bought successfully";
          // this.soldRecipes = append(this.soldRecipes, data.recipe);
          this.processResultStatus = true;
        } else {
          this.errorMessage = res.errors[0].message || "error buying recipe";
          this.processResultStatus = false;
        }
      } catch (error) {
        console.log("error buying a recipeeeee", error);
        this.errorMessage = error.message || "error buying recipe";
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
        return this.processResultStatus;
      }
    },
    async getCheckOutUrl(payload) {
      this.isLoading = true;
      const { $apollo } = useNuxtApp();
      try {
        const res = await $apollo.clients.default.query({
          query: CHECKOUT_URL,
          variables: {
            id: payload,
          },
        });
        if (res.data) {
          this.checkoutUrl = res.data.payments[0].checkout_url;
          console.log("checkout url", res);
          this.processResultStatus = true;
        }
      } catch (error) {
        console.error("error fetching checkout url", error);
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
      }
    },

    async deleteRecipes(payload) {
      this.isLoading = true;
      const { $apollo } = useNuxtApp();
      const { id, user_id } = payload;
      try {
        const res = await $apollo.clients.default.mutate({
          mutation: DELETE_RECIPE,
          variables: {
            id,
            user_id,
          },
        });
        if (res.data) {
          this.processResultStatus = true;
          this.recipes = this.recipes.filter((recipe) => recipe.id !== id);
        }
      } catch (error) {
        console.error("error deleting recipe", error);
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
      }
    },

    async getSellerRecipes() {
      this.isLoading = true;
      this.sellerId = Number(localStorage.getItem("authUserId"));
      const { $apollo } = useNuxtApp();

      try {
        const res = await $apollo.clients.default.query({
          query: SELLER_RECIPES,
          variables: {
            seller_id: this.sellerId,
          },
          fetchPolicy: "no-cache",
        });

        if (res.data) {
          this.soldRecipes = res.data.sold_recipes_aggregate.nodes;
          console.log(
            "Seller recipes:",
            JSON.stringify(this.soldRecipes, null, 2)
          );

          this.processResultStatus = true;
        }
      } catch (error) {
        console.error("Error fetching seller recipes:", error);
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
      }
    },

    async getPaymentDetails(payload) {
      const { recipe_id, buyer_id } = payload;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: PAYMENT_DETAILS,
          variables: {
            recipe_id,
            buyer_id,
          },
        });
        if (res) {
          const data = res.data?.sold_recipes?.[0]?.payments?.[0];
          console.log("payment detail", JSON.stringify(res, null, 2));
          console.log(
            "payment details",
            res.data?.sold_recipes?.[0]?.payments?.[0]?.payment_status
          );
          this.paymentDetails = data;
          this.tx_ref = data.tx_ref;
          this.paymentId = data.id;
          this.paymentStatus = data.payment_status;
        }
      } catch (error) {
        console.log("error fetching payment details", error);
      }
    },

    async verifyPayment(payload) {
      this.isLoading = true;
      const { $apollo } = useNuxtApp();
      const { id, tx_ref } = payload;

      try {
        const res = await $apollo.clients.default.mutate({
          mutation: VERIFY_PAYMENT,
          variables: {
            id,
            tx_ref,
          },
        });
        if (res.data) {
          console.log(
            JSON.stringify(
              "payment verification details",
              res.data.verifyPayment,
              null,
              2
            )
          );
          this.paymentStatus = res.data.verifyPayment.status;
          this.successMessage = res.data.verifyPayment.message;
          this.processResultStatus = true;
        } else {
          this.errorMessage = res.errors[0].message;
          this.processResultStatus = false;
        }
      } catch (error) {
        console.error("Error verifying payment", error);
      } finally {
        this.isLoading = false;
        return this.processResultStatus;
      }
    },
    async updateRecipe(payload) {
      this.isLoading = true;
      const {
        recipe_id,
        title,
        description,
        prep_time,
        cook_time,
        servings,
        category_id,
        price,
        user_id,
        imageName,
        imageType,
        base64string,
      } = payload;
      console.log("payload to update the recipe", payload);
      try {
        const { $apollo } = useNuxtApp();
        let variables = {
          recipe_id: recipe_id,
          user_id,
          title,
          description,
          prep_time,
          cook_time,
          servings,
          category_id,
          price,
          base64string: base64string,
          type: imageType,
          name: imageName,
        };

        console.log("recipe update variables", JSON.stringify(variables));
        const res = await $apollo.clients.default.mutate({
          mutation: UPDATE_RECIPE,
          variables,
          awaitRefetchQueries: true,
          refetchQueries: [
            {
              query: GET_RECIPE,
              variables: { recipeId: recipe_id },
            },
          ],
          onCompleted: (data) => {
            console.log("the refetched recipe data", data);
          },
          onError: (error) => {
            console.error("Error during mutation:", error);
          },
        });
        console.log("mutation response", res);

        if (res.data && res.data.updateRecipe) {
          const data = res.data.updateRecipe;
          this.updatedRecipe = data;
          console.log("result", JSON.stringify(this.updatedRecipe));
          this.successmessage = data.message || " recipe updated successfully!";
          this.processResultStatus = true;
        } else {
          throw new Error("Invalid respose from server");
        }
      } catch (error) {
        console.log(" error updating a recipes", error);
        this.errorMessage = "error updating the recipe";
        this.processResultStatus = false;
      } finally {
        this.isLoading = false;
        return this.processResultStatus;
      }
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useRecipeStore, import.meta.hot));
}