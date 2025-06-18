<script setup>
import { useRoute } from "vue-router";
import { ref } from "vue";
import { Form, Field } from "vee-validate";
import { useToast } from "vue-toast-notification";
import * as yup from "yup";

useSeoMeta({
  title: "recipe-app | Update User",
  description: "my cool recipe app",
});

const route = useRoute();
const router = useRouter();
const toast = useToast();

const userId = parseInt(route.params.id);

const getExtension = (filename) => {
  const extension = filename.split(".");
  return extension[1];
};
const MAX_FILE_SIZE = 2097152;
const useAuthStore = authStore();

const profileError = ref({
  userName: "",
  phone: "",
  profile: "",
});
const nameregex = /^([^\x00-\x7F]|[\a-zA-Z_\ \.\+\-]){2,20}$/;
const phoneregex = /^(^\+251|^251|^0)?(9|7)\d{8}$/;
const errorMessages = {
  userName: {
    invalid_username: "Invalid user name",
  },
  phone: {
    invalid_phone: "Invalid phone address",
  },
};

const schema = yup.object({
  userName: yup
    .string()
    .required()
    .matches(nameregex, errorMessages.userName.invalid_username),
  phone: yup
    .string()
    .required()
    .matches(phoneregex, errorMessages.phone.invalid_phone),
  profile: yup
    .mixed()
    .test({
      message: "Please provide supported file type",
      test: (file, context) => {
        let isValid;
        if (file) {
          isValid = [
            "jpg",
            "gif",
            "png",
            "jpeg",
            "svg",
            "webp",
            "jpeg",
          ].includes(getExtension(file?.name));
        } else {
          isValid = true;
        }
        if (!isValid) {
          context?.createError();
        }
        return isValid;
      },
    })
    .test({
      message: "File is too big, cant exit 2Mb",
      test: (file) => {
        let isValid;
        if (file) {
          isValid = file?.size < MAX_FILE_SIZE;
        } else {
          isValid = true;
        }
        return isValid;
      },
    }),
});

const handleFileChange = (event) => {
  const file = event.target.files[0];
  if (file) {
    image.value.name = file.name;
    image.value.type = file.type;
    const reader = new FileReader();
    reader.onload = () => {
      const base64String = reader.result.match(/base64,(.*)$/)[1];
      image.value.base64String = base64String;
    };
    reader.readAsDataURL(file);
    reader.onerror = (error) => {
      console.log("there is some error in image uploading: ", error);
      reject(error);
    };
  }
};

const image = ref({
  name: "",
  type: "",
  base64String: "",
});

const handleUpdateProfle = async (value) => {
  useAuthStore.changeOnLoad(true);
  let payload;
  if (
    image.value.name !== "" ||
    image.value.type !== "" ||
    image.value.base64String !== ""
  ) {
    payload = {
      userId: userId,
      userName: value.userName,
      phone: value.phone,
      imageName: image.value.name,
      imageType: image.value.type,
      imageString: image.value.base64String,
    };
  } else {
    payload = {
      userName: value.userName,
      phone: value.phone,
      userId: userId,
    };
  }
  const result = await useAuthStore.updateUser(payload);
  useAuthStore.changeOnLoad(false);
  if (result) {
    await useAuthStore.getProfile();
    if (
      useAuthStore.$state.successMessage &&
      useAuthStore.$state.successMessage.length > 0
    ) {
      const message = useAuthStore.$state.successMessage;
      useAuthStore.setSuccessMessage("");
      toast.success(message);
      router.push("/admin/dashboard");
    } else {
      toast.success("Your profile updated successfully!");
    }
  } else {
    if (
      useAuthStore.$state.errorMessage &&
      useAuthStore.$state.errorMessage.length > 0
    ) {
      const message = useAuthStore.$state.errorMessage;
      useAuthStore.setErrorMessage("");
      toast.error(message);
    } else {
      toast.error("Something wrong! please try again.");
    }
  }
};
</script>

<template>
  <div class="container mx-auto">
    <div class="mt-6 bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-bold mb-4 text-gray-800">Update User Profile</h2>
      <Form
        @submit="handleUpdateProfle"
        :validation-schema="schema"
        v-slot="{ errors }"
        class="space-y-4"
      >
        <div>
          <UIInput
            name="userName"
            type="text"
            placeholder="enter username"
            label="User Name"
            :is-required="true"
            :edit-mode="true"
            :value="
              useAuthStore.$state.user.username
                ? useAuthStore.$state.username
                : ''
            "
            :error-message="
              errors.userName
                ? errors.userName
                : profileError.userName
                ? profileError.userName
                : ''
            "
            class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:border-green-500 focus:ring-green-500"
          />
        </div>
        <div>
          <UIInput
            name="phone"
            placeholder="Insert your new phonenumber"
            label="Phone Number"
            :is-required="true"
            :edit-mode="true"
            :value="
              useAuthStore.$state.user.phone
                ? useAuthStore.$state.user.phone
                : ''
            "
            type="text"
            :error-message="
              errors.phone
                ? errors.phone
                : profileError.phone
                ? profileError.phone
                : ''
            "
          />
        </div>
        <div class="mt-6">
          <button
            class="w-full px-6 py-3 text-sm font-medium tracking-wide text-white capitalize transition-colors duration-300 transform bg-green-600 rounded-lg hover:bg-green-700 focus:outline-none focus:ring focus:ring-green-300 focus:ring-opacity-50"
          >
            <span v-if="useAuthStore.$state.onLoad == false">
              Update Profile
            </span>
            <div class="w-full flex items-center justify-center gap-2" v-else>
              <UILoading />
            </div>
          </button>
        </div>
      </Form>
    </div>
  </div>
</template>