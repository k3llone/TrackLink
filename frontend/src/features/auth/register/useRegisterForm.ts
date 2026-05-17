import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { registerUser } from "@/api/auth";
import type { ApiClientError, ApiFieldErrors } from "@/api/types";
import { useToast } from "@/shared/composables/useToast";
import { t } from "@/shared/lib/i18n";
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
      errors.email = t("auth.validation.emailRequired");
    } else if (!EMAIL_PATTERN.test(form.email.trim())) {
      errors.email = t("auth.validation.emailInvalid");
    }

    if (!form.password) {
      errors.password = t("auth.validation.passwordRequired");
    } else if (form.password.length < MIN_PASSWORD_LENGTH) {
      errors.password = t("auth.validation.passwordMin", { min: MIN_PASSWORD_LENGTH });
    }

    if (!form.confirmPassword) {
      errors.confirmPassword = t("auth.validation.confirmPasswordRequired");
    } else if (form.confirmPassword !== form.password) {
      errors.confirmPassword = t("auth.validation.passwordMismatch");
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

      toast.success(t("auth.register.successMessage"), t("auth.register.successTitle"));
      await router.push(ROUTES.login);

      return true;
    } catch (error: unknown) {
      if (isApiClientError(error)) {
        if (error.status === 409) {
          errors.email = t("auth.register.emailExists");
        } else if (error.status === 400 || error.status === 422) {
          applyFieldErrors(error.fields);
          errors.form = hasFieldErrors(error.fields) ? "" : t("auth.register.validationFailed");
        } else {
          errors.form = t("auth.register.submitFailed");
        }
      } else {
        errors.form = t("auth.register.networkFailed");
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
