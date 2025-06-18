import { acceptHMRUpdate, defineStore } from "pinia";
import USERS_QUERY from "../graphql/query/users.gql";
import GET_USER from "../graphql/query/user.gql";
import DELETE_USER from "../graphql/mutation/deleteUser.gql";
export const userStore = defineStore({
  id: "user",
  state: () => ({
    onload: false,
    users: [],
    errorMessage: "",
    successmessage: "",
    user: [],
    limit: 100,
    offset: 0,
    currentPage: 1,
    totalUsers: 0,
    searchUser: "%%",
    processResultStatus: false,
  }),
  actions: {
    setOffset(payload) {
      this.offset = payload;
    },
    setLimit(payload) {
      this.limit = payload;
    },
    setCurrentPage(payload) {
      if (payload < 1) {
        this.currentPage = 1;
      } else {
        this.currentPage = payload;
      }
      this.offset = (this.currentPage - 1) * this.limit;
    },
    setSearchUser(payload) {
      this.searchUser = "%" + payload + "%";
    },
    setErrorMessage(payload) {
      this.errorMessage = payload;
    },
    setSuccessMessage(payload) {
      this.successmessage = payload;
    },

    async getUsers(payload = "") {
      this.onload = true;
      this.setSearchUser(payload);

      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: USERS_QUERY,
          variables: {
            limit: this.limit,
            offset: this.offset,
          },
        });
        console.log("the user response data", JSON.stringify(res, null, 2));
        this.users = res.data.users;
        this.totalUsers = res.data.users.users_aggregate.aggregate.count;
      } catch (error) {
        console.log("error in loading user data", error);
      } finally {
        this.onload = false;
      }
    },

    async getUser(payload) {
      this.onload = true;
      try {
        const { $apollo } = useNuxtApp();
        const res = await $apollo.clients.default.query({
          query: GET_USER,
          variables: {
            userId: payload,
          },
        });
        console.log("the fetched user data", res);
        this.user = res.data.users_by_pk;
        console.log("the fetched user dataaaa", this.user);
      } catch (error) {
        console.log("error in loading the user data", error);
      } finally {
        this.onload = false;
      }
    },

    async deleteUser(payload) {
      this.onload = true;
      this.setSearchUser("");
      const { $apollo } = useNuxtApp();
      try {
        const res = await $apollo.clients.default.mutate({
          mutation: DELETE_USER,
          variables: {
            userId: payload,
          },
          awaitRefetchQueries: true,
          refetchQueries: [
            {
              query: USERS_QUERY,
              variables: {
                limit: this.limit,
                offset: this.offset,
              },
            },
          ],
          onCompleted: (data) => {
            console.log("the refetched user data", data);
          },
          onError: (error) => {
            console.error("Error during mutation:", error);
          },
        });
        if (res.data) {
          this.successmessage = "user deleted successfully!";
          this.processResultStatus = true;
        } else {
          this.errorMessage = res.errors[0].message;
          this.processResultStatus = false;
        }
      } catch (error) {
        console.log("error deleting the user", error);
        this.errorMessage = "unable to delete the user";
        this.processResultStatus = false;
      } finally {
        this.onload = false;
        return this.processResultStatus;
      }
    },
  },
});
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(userStore, import.meta.hot));
}