import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { createLink } from "@/api/links";
import type { ApiClientError, ApiFieldErrors } from "@/api/types";
import type { CreateLinkRequest, Link } from "@/entities/link/link.types";
import { useSession } from "@/entities/session/useSession";
import { useToast } from "@/shared/composables/useToast";
import { ROUTES } from "@/shared/lib/routes/paths";

const CUSTOM_ALIAS_PATTERN = /^[a-zA-Z0-9_-]+$/;
const MIN_CUSTOM_ALIAS_LENGTH = 3;
const MAX_CUSTOM_ALIAS_LENGTH = 64;

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const hasFieldErrors = (fields: ApiFieldErrors | null) => Boolean(fields && Object.keys(fields).length > 0);

const isValidTargetUrl = (value: string) => {
  try {
    const url = new URL(value);
    return (url.protocol === "http:" || url.protocol === "https:") && Boolean(url.hostname);
  } catch {
    return false;
  }
};

export const useCreateLinkForm = () => {
  const router = useRouter();
  const session = useSession();
  const toast = useToast();
  const form = reactive({
    targetUrl: "",
    customAlias: "",
  });
  const errors = reactive({
    targetUrl: "",
    customAlias: "",
    form: "",
  });
  const isSubmitting = ref(false);

  const clearErrors = () => {
    errors.targetUrl = "";
    errors.customAlias = "";
    errors.form = "";
  };

  const validate = () => {
    clearErrors();

    const targetUrl = form.targetUrl.trim();
    const customAlias = form.customAlias.trim();

    if (!targetUrl) {
      errors.targetUrl = "Укажите целевой URL";
    } else if (!isValidTargetUrl(targetUrl)) {
      errors.targetUrl = "Введите корректный http или https URL";
    }

    if (customAlias) {
      if (customAlias.length < MIN_CUSTOM_ALIAS_LENGTH || customAlias.length > MAX_CUSTOM_ALIAS_LENGTH) {
        errors.customAlias = "Alias должен содержать от 3 до 64 символов";
      } else if (!CUSTOM_ALIAS_PATTERN.test(customAlias)) {
        errors.customAlias = "Alias может содержать только латиницу, цифры, _ и -";
      }
    }

    return !errors.targetUrl && !errors.customAlias;
  };

  const getPayload = (): CreateLinkRequest => {
    const targetUrl = form.targetUrl.trim();
    const customAlias = form.customAlias.trim();
    const payload: CreateLinkRequest = { targetUrl };

    if (customAlias) {
      payload.customAlias = customAlias;
    }

    return payload;
  };

  const applyFieldErrors = (fields: ApiFieldErrors | null) => {
    if (!fields) {
      return;
    }

    errors.targetUrl = fields.targetUrl || fields.target_url || "";
    errors.customAlias = fields.customAlias || fields.custom_alias || "";
  };

  const handleUnauthorized = async () => {
    session.clearSession();
    await router.push(ROUTES.login);
  };

  const submit = async (): Promise<Link | null> => {
    if (isSubmitting.value) {
      return null;
    }

    if (!validate()) {
      return null;
    }

    isSubmitting.value = true;

    try {
      const link = await createLink(getPayload());
      toast.success("Короткая ссылка создана.");

      return link;
    } catch (error: unknown) {
      if (isApiClientError(error)) {
        if (error.status === 401) {
          await handleUnauthorized();
        } else if (error.status === 409) {
          errors.customAlias = "Этот alias уже занят";
        } else if (error.status === 400 || error.status === 422) {
          applyFieldErrors(error.fields);
          errors.form = hasFieldErrors(error.fields) ? "" : "Проверьте target URL и custom alias.";
        } else {
          errors.form = "Не удалось создать короткую ссылку. Повторите попытку позже.";
        }
      } else {
        errors.form = "Не удалось создать короткую ссылку. Проверьте соединение и попробуйте снова.";
      }

      return null;
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
