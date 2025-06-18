<script setup>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import { useToast } from "vue-toast-notification";

const toast = useToast();

useSeoMeta({
  title: "recipe-app | Create User",
  description: "my cool recipe app",
});
definePageMeta({
  layout: "admin",
});

const useAuthStore = authStore();
const MAX_FILE_SIZE = 2097152;
const router = useRouter();

const getExtension = (filename) => {
  const extension = filename.split(".");
  return extension[1];
};

const image = ref({
  name: "",
  type: "",
  base64String: "",
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

const signupError = ref({
  userName: "",
  email: "",
  password: "",
  passwordConfirmation: "",
  phone: "",
  profile: "",
});

const nameregex = /^([^\x00-\x7F]|[\a-zA-Z_\ \.\+\-]){2,20}$/;
const emailregex = /^((?!\.)[\w_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$/gm;
const phoneregex = /^(^\+251|^251|^0)?(9|7)\d{8}$/;
const passwordCapitalLetter = /^((?=\S*?[A-Z]).{8,})\S$/;
const passwordNumber = /^((?=\S*?[0-9]).{8,})\S$/;

const errorMessages = {
  email: {
    required: "Email field is required!",
    not_email: "Sorry it is not an email address!",
    invalid_email: "Invalid email address!",
  },
  password: {
    required: "Password field is required!",
    min_password_len: "The minimum length of password is 6",
    not_valid_password_capitalizer:
      "Atleast there is 1 capital letter in the password",
    not_vaild_password_digit: "Password also must have a digit in it",
  },
  userName: {
    invalid_username: "Invalid user name",
  },
  passwordConfirmation: {
    must_be_same: "Password confirmation must be same with password",
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
  email: yup
    .string()
    .required(errorMessages.email.required)
    .email(errorMessages.email.not_email)
    .matches(emailregex, errorMessages.email.invalid_email),
  password: yup
    .string()
    .required(errorMessages.password.required_password)
    .min(8, errorMessages.password.min_password_len)
    .matches(
      passwordCapitalLetter,
      errorMessages.password.not_valid_password_capitalzer
    )
    .matches(passwordNumber, errorMessages.password.not_vaild_password_digit),
  passwordConfirmation: yup
    .string()
    .required()
    .oneOf(
      [yup.ref("password"), null],
      errorMessages.passwordConfirmation.must_be_same
    ),
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

const handleSigningUp = async (value) => {
  console.log("Form Data from signup:", value);
  useAuthStore.changeOnLoad(true);
  let payload;
  if (
    image.value.name !== "" ||
    image.value.type !== "" ||
    image.value.base64String !== ""
  ) {
    payload = {
      userName: value.userName,
      email: value.email,
      password: value.password,
      phone: value.phone,
      imageName: image.value.name,
      imageType: image.value.type,
      imageString: image.value.base64String,
    };
  } else {
    payload = {
      userName: value.userName,
      email: value.email,
      password: value.password,
      phone: value.phone,
    };
  }
  const result = await useAuthStore.signupUser(payload);
  router.push("/admin/dashboard");
  if (result) {
    toast.success(
      "Welcome to my app, check your email to verify your account!"
    );
    router.push("/admin/dashboard");
  } else {
    if (useAuthStore.$state.errorMessage) {
      const message = useAuthStore.$state.errorMessage;
      // toast.error(message);
      useAuthStore.setErrorMessage("");
    } else {
      toast.error("someting went wrong");
    }
  }
};
</script>

<template>
  <div
    class="container mx-auto mt-10 w-full border-2 border-green-100 shadow-lg rounded-lg items-center justify-center px-2 sm:px-6 lg:flex-none lg:px-4 xl:px-6"
  >
    <h1 class="text-4xl font-bold items-center text-gray-800">Create User</h1>
    <div class="w-full pt-2 px-8 py-8 md:px-8">
      <Form
        @submit="handleSigningUp"
        :validation-schema="schema"
        v-slot="{ errors }"
        class="w-full flex-col gap-4 items-center justify-center"
      >
        <UIInput
          name="userName"
          placeholder="Insert Your Full Name"
          label="User Name"
          :error-message="
            errors.userName
              ? errors.userName
              : signupError.userName
              ? signupError.userName
              : ''
          "
          :is-required="true"
        />
        <UIInput
          name="email"
          placeholder="Insert Your Email"
          label="Email"
          :error-message="
            errors.email
              ? errors.email
              : signupError.email
              ? signupError.email
              : ''
          "
          :is-required="true"
        />
        <UIInput
          name="password"
          placeholder="Insert Your Password"
          label="password"
          :is-required="true"
          type="password"
          :error-message="
            errors.password
              ? errors.password
              : signupError.password
              ? signupError.password
              : ''
          "
          :is-password="true"
        />
        <UIInput
          name="passwordConfirmation"
          placeholder="Confirm Your Password"
          label="Password Confirmation"
          :is-required="true"
          type="password"
          :error-message="
            errors.passwordConfirmation
              ? errors.passwordConfirmation
              : signupError.passwordConfirmation
              ? signupError.passwordConfirmation
              : ''
          "
          :is-password="true"
        />
        <UIInput
          name="phone"
          placeholder="Insert your phone"
          label="Phonenumber"
          :is-required="true"
          type="text"
          :error-message="
            errors.phone
              ? errors.phone
              : signupError.phone
              ? signupError.phone
              : ''
          "
        />
        <div class="w-full flex flex-col sm:mb-3 justify-center gap-2">
          <label for="" class="pl-3 ml-px text-sm font-medium text-gray-700"
            >Profile
          </label>
          <UIImage
            name="profile"
            @image-changed="handleFileChange"
            :error-message="
              errors.profile
                ? errors.profile
                : signupError.profile
                ? signupError.profile
                : ''
            "
          />
        </div>
        <div class="flex justify-center mt-4">
          <button
            class="bg-green-500 text-white rounded-lg px-6 mb-2 py-3 hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50 transition-colors duration-300"
          >
            <span v-if="useAuthStore.$state.onLoad == false">Create User</span>
            <div class="w-full flex items-center justify-center gap-2" v-else>
              <UILoading />
            </div>
          </button>
        </div>
      </Form>
    </div>
  </div>
</template>