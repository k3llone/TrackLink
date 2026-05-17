import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { logoutUser } from "@/api/auth";
import type { ApiClientError } from "@/api/types";
import { useSession } from "@/entities/session/useSession";
import { useToast } from "@/shared/composables/useToast";
import { t } from "@/shared/lib/i18n";
import { ROUTES } from "@/shared/lib/routes/paths";

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

export const useLogout = () => {
  const router = useRouter();
  const session = useSession();
  const toast = useToast();
  const isLoggingOut = ref(false);

  const clearAndRedirectToLogin = async () => {
    session.clearSession();
    await router.push(ROUTES.login);
  };

  const logout = async () => {
    if (isLoggingOut.value) {
      return false;
    }

    isLoggingOut.value = true;

    try {
      await logoutUser();
      await clearAndRedirectToLogin();

      return true;
    } catch (error: unknown) {
      if (isApiClientError(error) && error.status === 401) {
        await clearAndRedirectToLogin();
        return true;
      }

      toast.error(t("session.networkError"));
      return false;
    } finally {
      isLoggingOut.value = false;
    }
  };

  return {
    isLoggingOut: computed(() => isLoggingOut.value),
    logout,
  };
};
