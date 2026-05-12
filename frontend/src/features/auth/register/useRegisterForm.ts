import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { registerUser } from "@/api/auth";
import type { ApiClientError, ApiFieldErrors } from "@/api/types";
import { useToast } from "@/shared/composables/useToast";
import { ROUTES } from "@/shared/lib/routes/paths";

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const MIN_PASSWORD_LENGTH = 8;

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const hasFieldErrors = (fields: ApiFieldErrors | null) => Boolean(fields && Object.keys(fields).length > 0);

export const useRegisterForm = () => {
  const router = useRouter();
  const toast = useToast();
  const form = reactive({
    email: "",
    password: "",
    confirmPassword: "",
  });
  const errors = reactive({
    email: "",
    password: "",
    confirmPassword: "",
    form: "",
  });
  const isSubmitting = ref(false);

  const clearErrors = () => {
    errors.email = "";
    errors.password = "";
    errors.confirmPassword = "";
    errors.form = "";
  };

  const validate = () => {
    clearErrors();

    if (!form.email.trim()) {
      errors.email = "Email is required";
    } else if (!EMAIL_PATTERN.test(form.email.trim())) {
      errors.email = "Enter a valid email address";
    }

    if (!form.password) {
      errors.password = "Password is required";
    } else if (form.password.length < MIN_PASSWORD_LENGTH) {
      errors.password = `Password must be at least ${MIN_PASSWORD_LENGTH} characters`;
    }

    if (!form.confirmPassword) {
      errors.confirmPassword = "Confirm your password";
    } else if (form.confirmPassword !== form.password) {
      errors.confirmPassword = "Passwords do not match";
    }

    return !errors.email && !errors.password && !errors.confirmPassword;
  };

  const applyFieldErrors = (fields: ApiFieldErrors | null) => {
    if (!fields) {
      return;
    }

    errors.email = fields.email || "";
    errors.password = fields.password || "";
  };

  const submit = async () => {
    if (!validate()) {
      return false;
    }

    isSubmitting.value = true;

    try {
      await registerUser({
        email: form.email.trim(),
        password: form.password,
      });

      toast.success("Account created. Please sign in.", "Registration complete");
      await router.push(ROUTES.login);

      return true;
    } catch (error: unknown) {
      if (isApiClientError(error)) {
        if (error.status === 409) {
          errors.email = "Email already exists";
        } else if (error.status === 400 || error.status === 422) {
          applyFieldErrors(error.fields);
          errors.form = hasFieldErrors(error.fields) ? "" : "Проверьте корректность email и пароля.";
        } else {
          errors.form = "Не удалось создать аккаунт. Попробуйте еще раз позже.";
        }
      } else {
        errors.form = "Не удалось создать аккаунт. Проверьте соединение и попробуйте снова.";
      }

      return false;
    } finally {
      isSubmitting.value = false;
    }
  };

  return {
    form,
    errors,
    isSubmitting: computed(() => isSubmitting.value),
    submit,
  };
};
