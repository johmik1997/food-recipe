import { acceptHMRUpdate, defineStore } from "pinia";
import SIGNUP_MUTATION from "../graphql/mutation/signup.gql";
import VERIFY_USER from "../graphql/mutation/verifyEmail.gql";
import LOGIN_USER from "../graphql/query/login.gql";
import FORGOT_PASS from "../graphql/mutation/forgotPassword.gql";
import RESET_PASS from "../graphql/mutation/passwordReset.gql";
import UPDATE_PROFILE from "../graphql/mutation/updateProfile.gql";
import PROFILE_QUERY from "../graphql/query/profile.gql";

export const authStore = defineStore(
  {
    id: "auth",
    state: () => ({
      onLoad: false,
      user: {},
      profileLoading: false,
      isAuthenticated: false,
      isAdmin: false,
      isAuthed: false,
      authPages: ["AuthPage"],
      anonymousPages: ["index", "Password-reset", "Verify", "Welcome", "Auth"],
      protectedPages: ["profile"],
      adminPages: ["admin-dashboard", "admin-users-id"],
      role: "",
      authToken: "",
      userId: "",
      userName: "",
      errorMessage: "",
      successMessage: "",
      processResultStatus: false,
      isEmailVerified: false,
    }),
    actions: {
      initAuthState() {
        this.isAuthed = !!localStorage.getItem("authToken");
        this.authToken = localStorage.getItem("authToken") || "";
        this.role = localStorage.getItem("authRole") || "anonymous";
        this.isAdmin = this.role == "systemAdmin";
        this.userId = localStorage.getItem("authUserId") || "";
      },

      changeOnLoad(payload) {
        this.onLoad = payload;
      },

      setProfileLoading(payload) {
        this.profileLoading = payload;
      },
      setErrorMessage(payload) {
        this.errorMessage = payload;
      },
      setUserId(payload) {
        this.userId = payload;
      },

      setSuccessMessage(payload) {
        this.successMessage = payload;
      },
      setProcessResultStatus(payload) {
        this.processResultStatus = payload;
      },
      setIsEmailVerified(payload) {
        this.isEmailVerified = payload;
      },
      setNotifications(payload) {
        this.notifications = payload;
      },
      setUser(payload) {
        this.user = payload;
      },

      isAdminPage(page) {
        return this.adminPages.includes(page);
      },

      isProtectedPage(page) {
        return this.protectedPages.includes(page);
      },

      isAnonymousPage(page) {
        console.log(`Anonymous pages: ${this.anonymousPages}`);
        return this.anonymousPages.includes(page);
      },

      isAuthPage(page) {
        return this.authPages.includes(page);
      },
      getRole() {
        const role = localStorage.getItem("authRole");
        if (role) {
          this.$state.role = role;
          this.$state.isAdmin = role === "systemAdmin";
          console.log("Role set in state:", this.$state.role);
        } else {
          this.$state.role = "user";
          this.$state.isAdmin = false;
          console.log("No role found in localStorage, setting default role.");
        }
      },

      async verifyEmail(payload) {
        console.log("verifyEmail function called with payload:", payload);
        const { user_id, verification_token } = payload;
        const { $apollo } = useNuxtApp();

        try {
          this.onLoad = true;
          const res = await $apollo.clients.default.mutate({
            mutation: VERIFY_USER,
            variables: {
              user_id,
              verification_token,
            },
          });
          console.log("Mutation response from the store:", res);
          if (res.data) {
            const data = res.data.verifyEmail.message;
            console.log("verification sucess:", data);
            this.isEmailVerified = true;
            this.successMessage = "Email verified sucessfully!";
          } else {
            this.errorMessage =
              res.errors?.[0]?.message || "verification failed";
          }
        } catch (error) {
          console.error("verififcation errror", error);
          this.errorMessage =
            error.message || "An error occuers during verification";
        } finally {
          this.onLoad = false;
        }
      },

      async loginUser(payload) {
        const { email, password } = payload;
        const { $apollo } = useNuxtApp();

        try {
          const res = await $apollo.clients.default.query({
            query: LOGIN_USER,
            variables: {
              email,  
              password,
            },
          });
          if (res.data) {
            const data = res.data.login.user;
            console.log("login data", data);
            this.authToken = data.token;
            this.userId = data.id;
            this.userName = data.name;
            this.isAuthed = true;
            this.user.email = data.email;
            this.user.userName = data.name;
            this.user.role = data.role;
            this.role = data.role;
            localStorage.setItem("authRole", this.role);
            if (data.role == "systemAdmin") {
              this.isAdmin = true;
            }
            localStorage.setItem("authUserId", data.id);
            localStorage.setItem("authToken", data.token);
            this.processResultStatus = true;
          } else {
            this.errorMessage = res.errors[0].message;
            this.processResultStatus = false;
          }
        } catch (error) {
          console.log("login error", error);
          if (error.message) {
            this.setErrorMessage(error.message);
          }
          this.processResultStatus = false;
        } finally {
          this.onLoad = false;
          return this.processResultStatus;
        }
      },

      async logout(payload) {
        localStorage.removeItem("authToken");
        localStorage.removeItem("authRole");
        localStorage.removeItem("authUserId");
        this.authToken = "";
        this.userId = "";
        this.user = {};
        this.isAdmin = false;
        this.isAuthed = false;
        this.role = "";
        return true;
      },

      async forgotPassword(payload) {
        const { email } = payload;
        const { $apollo } = useNuxtApp();

        try {
          const response = await $apollo.clients.default.mutate({
            mutation: FORGOT_PASS,
            variables: { email },
          });
          if (response.data) {
            const data = response.data;
            console.log("the forgot pass data", data);
            this.processResultStatus = true;
          } else {
            this.errorMessage = response.errors[0].message;
            this.processResultStatus = false;
          }
        } catch (error) {
          console.log("forgot password error");
          if (error.message) {
            this.errorMessage = error.message;
          }
          this.processResultStatus = false;
        } finally {
          this.onLoad = false;
          return this.processResultStatus;
        }
      },
      async resetPassword(payload) {
        const { password, token, userId } = payload;
        const { $apollo } = useNuxtApp();
        try {
          this.onLoad = true;
          const response = await $apollo.clients.default.mutate({
            mutation: RESET_PASS,
            variables: { password, token, userId },
          });
          if (response.data) {
            const data = response.data;
            console.log("the reset password data:", data);
            if (data.passwordUpdate) {
              this.successMessage = data.passwordUpdate.message
                ? data.passwordUpdate.message
                : "";
            }
            this.onLoad = false;
            this.processResaltStatus = true;
          } else {
            this.errorMessage = response.errors[0].message;
            this.processResaltStatus = false;
          }
        } catch (error) {
          console.log("reset password error:", error);
          if (error.message) {
            this.errorMessage = error.message;
          }
          this.processResaltStatus = false;
        } finally {
          this.onLoad = false;
          return this.processResaltStatus;
        }
      },

      async updateUser(payload) {
        const { userId, userName, phone, imageName, imageType, imageString } =
          payload;

        const { $apollo } = useNuxtApp();

        try {
          let variables;
          const userIdToUpdate =
            userId ?? Number(localStorage.getItem("authUserId"));

          // const userId = localStorage.getItem("authUserId");
          const userIdInt = Number(userId);

          if (imageName || imageType || imageString) {
            variables = {
              userId: userIdToUpdate,
              userName,
              phone,
              base64String: imageString,
              imageType: imageType,
              imageName: imageName,
            };
          } else {
            variables = {
              userId: userIdToUpdate,
              userName,
              phone,
            };
          }

          this.userId = userIdInt;
          const res = await $apollo.clients.default.mutate({
            mutation: UPDATE_PROFILE,
            variables,
            awaitRefetchQueries: true,
            refetchQueries: [
              {
                query: PROFILE_QUERY,
                variables: { id: this.userId },
              },
            ],
            onCompleted: (data) => {
              console.log("the fetched user data", data);
            },

            onError: (error) => {
              console.log("error during fething", error);
            },
          });
          if (res.data) {
            console.log("update profile responseeee", res.data);
            if (res.data.updateProfile.message) {
              this.successMessage = res.data.updateProfile.message;
              this.processResultStatus = true;
              this.user = res.data.updateProfile;
            } else {
              this.successMessage = "your profile has been updated succesfully";
              this.processResultStatus = true;
            }
          } else {
            this.errorMessage = res.errors[0].message;
            this.processResultStatus = false;
          }
        } catch (error) {
          console.log("update profile error", error);
          if (error.message) {
            this.setErrorMessage(error.message);
          }
          this.processResultStatus = false;
        } finally {
          this.onLoad = false;
          return this.processResultStatus;
        }
      },

      async updateProfile(payload) {
        const { userName, phone, imageName, imageType, imageString } = payload;

        const { $apollo } = useNuxtApp();

        try {
          let variables;
          const userId = localStorage.getItem("authUserId");
          const userIdInt = Number(userId);

          if (imageName || imageType || imageString) {
            variables = {
              userName,
              phone,
              base64String: imageString,
              imageType: imageType,
              imageName: imageName,
              userId: userIdInt,
            };
          } else {
            variables = {
              userName,
              phone,
              userId: userIdInt,
            };
          }

          this.userId = userIdInt;
          const res = await $apollo.clients.default.mutate({
            mutation: UPDATE_PROFILE,
            variables,
            awaitRefetchQueries: true,
            refetchQueries: [
              {
                query: PROFILE_QUERY,
                variables: { id: this.userId },
              },
            ],
            onCompleted: (data) => {
              console.log("the fetched user data", data);
            },

            onError: (error) => {
              console.log("error during fething", error);
            },
          });
          if (res.data) {
            console.log("update profile responseeee", res.data);
            if (res.data.updateProfile.message) {
              this.successMessage = res.data.updateProfile.message;
              this.processResultStatus = true;
              this.user = res.data.updateProfile;
            } else {
              this.successMessage = "your profile has been updated succesfully";
              this.processResultStatus = true;
            }
          } else {
            this.errorMessage = res.errors[0].message;
            this.processResultStatus = false;
          }
        } catch (error) {
          console.log("update profile error", error);
          if (error.message) {
            this.setErrorMessage(error.message);
          }
          this.processResultStatus = false;
        } finally {
          this.onLoad = false;
          return this.processResultStatus;
        }
      },

      async getProfile(payload) {
        this.profileLoading = true;
        this.userId = Number(localStorage.getItem("authUserId"));
        console.log("User ID from store:", this.userId);

        if (!isNaN(this.userId)) {
          const { $apollo } = useNuxtApp();
          try {
            const res = await $apollo.clients.default.query({
              query: PROFILE_QUERY,
              variables: { id: this.userId },
            });

            if (res.data && res.data.users_by_pk) {
              console.log("User data:", JSON.stringify(res, null, 2));
              const user = res.data.users_by_pk;
              this.user = user;

              if (user) {
                this.notifications = user.notifications;

                localStorage.setItem("authRole", user.role);
                this.isAdmin = user.role === "systemAdmin";
              }
              this.processResultStatus = true;
            } else {
              this.errorMessage =
                res.errors?.[0]?.message || "Unknown error occurred.";
              this.processResultStatus = false;
            }
          } catch (error) {
            console.error("Error getting user data:", error);
            this.errorMessage = error.message || "Failed to fetch user data.";
            this.processResultStatus = false;
          } finally {
            this.profileLoading = false;
          }
        } else {
          console.error("Invalid user ID.");
          this.errorMessage = "Invalid user ID.";
          this.processResultStatus = false;
        }

        return this.processResultStatus;
      },

      async signupUser(payload) {
        const {
          email,
          password,
          userName,
          phone,
          imageName,
          imageType,
          imageString,
        } = payload;

        const { $apollo } = useNuxtApp();

        try {
          let variables;
          if (imageName || imageType || imageString) {
            variables = {
              email,
              password,
              userName,
              phone,
              base64String: imageString,
              type: imageType,
              name: imageName,
            };
          } else {
            variables = {
              email,
              password,
              userName,
              phone,
            };
          }
          const res = await $apollo.clients.default.mutate({
            mutation: SIGNUP_MUTATION,
            variables,
          });
          console.log(res);

          if (res.data) {
            const data = res.data.signup;
            console.log("the signup data:", data);
            this.authToken = data.token;
            this.userId = data.id;
            this.isAuthed = true;
            console.log("is authed", this.isAuthed);
            this.user.email = data.email;
            this.user.userName = data.userName;
            console.log("name", this.user.userName);
            this.user.role = data.role;
            this.role = data.role;
            localStorage.setItem("authToken", data.token);
            localStorage.setItem("authRole", data.role);
            if (data.role == "systemAdmin") {
              this.isAdmin = true;
            }
            localStorage.setItem("authUserId", data.id);
            this.processResultStatus = true;
          } else {
            this.errorMessage = res.errors[0].message;
            this.processResultStatus = false;
          }
        } catch (error) {
          console.log("signup error", error);
          if (error.message) {
            this.setErrorMessage(error.message);
          }
          this.processResultStatus = false;
        } finally {
          this.onLoad = false;
          return this.processResultStatus;
        }
      },
    },
  }
  // {}
);
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(authStore, import.meta.hot));
}
