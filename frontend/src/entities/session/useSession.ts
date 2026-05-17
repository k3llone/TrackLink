import { computed, readonly, ref } from "vue";
import { getCurrentUser } from "@/api/auth";
import type { ApiClientError } from "@/api/types";
import { t } from "@/shared/lib/i18n";
import type { SessionError, SessionStatus, User } from "./session.types";

const user = ref<User | null>(null);
const status = ref<SessionStatus>("idle");
const sessionError = ref<SessionError | null>(null);
let pendingUserRequest: Promise<User | null> | null = null;

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const applyUser = (nextUser: User | null) => {
  user.value = nextUser;
  status.value = nextUser ? "authenticated" : "guest";
  sessionError.value = null;

  return nextUser;
};

const normalizeSessionError = (error: unknown): SessionError => {
  if (isApiClientError(error)) {
    if (error.status === 401) {
      return {
        status: error.status,
        code: "unauthorized",
        message: t("session.expired"),
      };
    }

    if (error.status === 403) {
      return {
        status: error.status,
        code: "forbidden",
        message: t("session.forbidden"),
      };
    }

    return {
      status: error.status,
      code: error.code || "session_check_failed",
      message: t("session.verifyFailed"),
    };
  }

  return {
    status: null,
    code: "network_error",
    message: t("session.networkError"),
  };
};

const applySessionError = (error: unknown) => {
  const nextError = normalizeSessionError(error);

  user.value = null;
  status.value = "guest";
  sessionError.value = nextError.status === 401 ? null : nextError;

  return null;
};

const loadCurrentUser = async (options: { force?: boolean } = {}) => {
  if (!options.force && status.value === "authenticated") {
    return user.value;
  }

  if (!options.force && status.value === "guest" && !sessionError.value) {
    return user.value;
  }

  if (pendingUserRequest) {
    return pendingUserRequest;
  }

  status.value = "loading";
  sessionError.value = null;
  pendingUserRequest = getCurrentUser()
    .then((currentUser) => applyUser(currentUser))
    .catch((error: unknown) => applySessionError(error))
    .finally(() => {
      pendingUserRequest = null;
    });

  return pendingUserRequest;
};

const setUser = (nextUser: User) => applyUser(nextUser);

const clearSession = () => {
  applyUser(null);
};

export const useSession = () => ({
  user: readonly(user),
  status: readonly(status),
  sessionError: readonly(sessionError),
  isLoading: computed(() => status.value === "loading"),
  isSessionLoading: computed(() => status.value === "loading"),
  isAuthenticated: computed(() => status.value === "authenticated"),
  isAdmin: computed(() => user.value?.role === "admin"),
  loadCurrentUser,
  setUser,
  clearSession,
});
