<script setup lang="ts">
import { computed, ref } from "vue";
import { deactivateAdminLink, type AdminLink } from "@/api/admin";
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
  deactivated: [link: AdminLink];
}>();

const toast = useToast();
const isDeactivating = ref(false);

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const canDeactivate = computed(() => props.link.status === "active");
const isDisabled = computed(() => props.disabled || isDeactivating.value || !canDeactivate.value);
const title = computed(() =>
  canDeactivate.value ? "Деактивировать ссылку" : "Можно деактивировать только активную ссылку.",
);

const getDeactivateLinkErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401 || error.status === 403) {
      return "У вас нет доступа к административной панели.";
    }

    if (error.status === 404) {
      return "Ссылка не найдена или уже удалена.";
    }

    if (error.status === 409 || error.code === "status_change_not_allowed") {
      return "Можно деактивировать только активную ссылку.";
    }
  }

  return "Не удалось деактивировать ссылку.";
};

const deactivate = async () => {
  if (isDisabled.value) {
    return;
  }

  isDeactivating.value = true;

  try {
    const updatedLink = await deactivateAdminLink(props.link.id);
    toast.success("Ссылка деактивирована.");
    emit("deactivated", updatedLink);
  } catch (error: unknown) {
    toast.error(getDeactivateLinkErrorMessage(error));
  } finally {
    isDeactivating.value = false;
  }
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isDeactivating"
    :full-width="fullWidth"
    :title="title"
    @click.stop="deactivate"
  >
    Деактивировать
  </UiButton>
</template>
