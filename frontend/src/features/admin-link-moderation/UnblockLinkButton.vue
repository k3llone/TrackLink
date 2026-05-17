<script setup lang="ts">
import { computed, ref } from "vue";
import { unblockAdminLink, type AdminLink } from "@/api/admin";
import type { ApiClientError } from "@/api/types";
import { useI18n } from "@/shared/composables/useI18n";
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
const { t } = useI18n();
const isUnblocking = ref(false);

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const canUnblock = computed(() => props.link.status === "blocked");
const isDisabled = computed(() => props.disabled || isUnblocking.value || !canUnblock.value);
const title = computed(() => (canUnblock.value ? t("admin.unblock.title") : t("admin.unblock.disabledTitle")));

const getUnblockLinkErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401 || error.status === 403) {
      return t("admin.unblock.error.access");
    }

    if (error.status === 404) {
      return t("admin.unblock.error.notFound");
    }
  }

  return t("admin.unblock.error.failed");
};

const unblock = async () => {
  if (isDisabled.value) {
    return;
  }

  isUnblocking.value = true;

  try {
    const updatedLink = await unblockAdminLink(props.link.id);
    toast.success(t("admin.unblock.success"));
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
    {{ t("admin.unblock.label") }}
  </UiButton>
</template>
