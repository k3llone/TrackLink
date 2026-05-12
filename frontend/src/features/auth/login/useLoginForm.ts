import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { loginUser } from "@/api/auth";
import type { ApiClientError, ApiFieldErrors } from "@/api/types";
import { useSession } from "@/entities/session/useSession";
import { ROUTES } from "@/shared/lib/routes/paths";

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const hasFieldErrors = (fields: ApiFieldErrors | null) => Boolean(fields && Object.keys(fields).length > 0);

export const useLoginForm = () => {
  const router = useRouter();
  const session = useSession();
  const form = reactive({
    email: "",
    password: "",
  });
  const errors = reactive({
    email: "",
    password: "",
    form: "",
  });
  const isSubmitting = ref(false);

  const clearErrors = () => {
    errors.email = "";
    errors.password = "";
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
    }

    return !errors.email && !errors.password;
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
      const response = await loginUser({
        email: form.email.trim(),
        password: form.password,
      });

      session.setUser(response.user);
      await router.push(ROUTES.dashboard);

      return true;
    } catch (error: unknown) {
      if (isApiClientError(error)) {
        if (error.status === 401) {
          errors.form = "Invalid username or password";
        } else if (error.status === 400 || error.status === 422) {
          applyFieldErrors(error.fields);
          errors.form = hasFieldErrors(error.fields) ? "" : "Проверьте корректность email и пароля.";
        } else {
          errors.form = "Не удалось войти. Попробуйте еще раз позже.";
        }
      } else {
        errors.form = "Не удалось войти. Проверьте соединение и попробуйте снова.";
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
