import { computed, readonly, ref } from "vue";
import { getCurrentUser, type User } from "@/api/auth";
import type { ApiClientError } from "@/api/types";

type SessionStatus = "idle" | "loading" | "authenticated" | "guest";

const user = ref<User | null>(null);
const status = ref<SessionStatus>("idle");
let pendingUserRequest: Promise<User | null> | null = null;

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const applyUser = (nextUser: User | null) => {
  user.value = nextUser;
  status.value = nextUser ? "authenticated" : "guest";

  return nextUser;
};

const loadCurrentUser = async (options: { force?: boolean } = {}) => {
  if (!options.force && status.value !== "idle") {
    return user.value;
  }

  if (pendingUserRequest) {
    return pendingUserRequest;
  }

  status.value = "loading";
  pendingUserRequest = getCurrentUser()
    .then((currentUser) => applyUser(currentUser))
    .catch((error: unknown) => {
      applyUser(null);

      if (isApiClientError(error) && error.status === 401) {
        return null;
      }

      return null;
    })
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
  isLoading: computed(() => status.value === "loading"),
  isAuthenticated: computed(() => status.value === "authenticated"),
  isAdmin: computed(() => user.value?.role === "admin"),
  loadCurrentUser,
  setUser,
  clearSession,
});
