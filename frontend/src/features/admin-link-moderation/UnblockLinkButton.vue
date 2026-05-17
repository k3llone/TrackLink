<script setup lang="ts">
import { computed, ref } from "vue";
import { unblockAdminLink, type AdminLink } from "@/api/admin";
import type { ApiClientError } from "@/api/types";
import { useToast } from "@/shared/composables/useToast";
import { UiButton } from "@/shared/ui";

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

const props = withDefaults(
  defineProps<{
    link: AdminLink;
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
    fullWidth?: boolean;
  }>(),
  {
    variant: "ghost",
    size: "sm",
    disabled: false,
    fullWidth: false,
  },
);

const emit = defineEmits<{
  unblocked: [link: AdminLink];
}>();

const toast = useToast();
const isUnblocking = ref(false);

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const canUnblock = computed(() => props.link.status === "blocked");
const isDisabled = computed(() => props.disabled || isUnblocking.value || !canUnblock.value);
const title = computed(() =>
  canUnblock.value ? "Разблокировать ссылку" : "Можно разблокировать только заблокированную ссылку.",
);

const getUnblockLinkErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401 || error.status === 403) {
      return "У вас нет доступа к административной панели.";
    }

    if (error.status === 404) {
      return "Ссылка не найдена или уже удалена.";
    }
  }

  return "Не удалось разблокировать ссылку.";
};

const unblock = async () => {
  if (isDisabled.value) {
    return;
  }

  isUnblocking.value = true;

  try {
    const updatedLink = await unblockAdminLink(props.link.id);
    toast.success("Ссылка разблокирована.");
    emit("unblocked", updatedLink);
  } catch (error: unknown) {
    toast.error(getUnblockLinkErrorMessage(error));
  } finally {
    isUnblocking.value = false;
  }
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isUnblocking"
    :full-width="fullWidth"
    :title="title"
    @click.stop="unblock"
  >
    Разблокировать
  </UiButton>
</template>
