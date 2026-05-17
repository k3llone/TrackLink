<script setup lang="ts">
import { computed, ref } from "vue";
import { blockAdminLink, type AdminLink } from "@/api/admin";
import type { ApiClientError } from "@/api/types";
import { useI18n } from "@/shared/composables/useI18n";
import { useToast } from "@/shared/composables/useToast";
import { UiButton } from "@/shared/ui";
import BlockLinkDialog from "./BlockLinkDialog.vue";

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

const props = withDefaults(
  defineProps<{
    link: AdminLink;
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
    fullWidth?: boolean;
    label?: string;
  }>(),
  {
    variant: "ghost",
    size: "sm",
    disabled: false,
    fullWidth: false,
    label: "",
  },
);

const emit = defineEmits<{
  blocked: [link: AdminLink];
}>();

const toast = useToast();
const { t } = useI18n();
const isDialogOpen = ref(false);
const isBlocking = ref(false);

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const canBlock = computed(() => props.link.status !== "blocked" && props.link.status !== "deleted");
const isDisabled = computed(() => props.disabled || isBlocking.value || !canBlock.value);
const buttonLabel = computed(() => props.label || t("admin.block.label"));
const title = computed(() => (canBlock.value ? buttonLabel.value : t("admin.block.disabledTitle")));

const getBlockLinkErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401 || error.status === 403) {
      return t("admin.block.error.access");
    }

    if (error.status === 404 || error.status === 409 || error.code === "status_change_not_allowed") {
      return t("admin.block.error.unavailable");
    }
  }

  return t("admin.block.error.failed");
};

const requestBlock = () => {
  if (isDisabled.value) {
    return;
  }

  isDialogOpen.value = true;
};

const confirmBlock = async () => {
  if (isDisabled.value) {
    return;
  }

  isBlocking.value = true;

  try {
    const updatedLink = await blockAdminLink(props.link.id);
    isDialogOpen.value = false;
    toast.success(t("admin.block.success"));
    emit("blocked", updatedLink);
  } catch (error: unknown) {
    toast.error(getBlockLinkErrorMessage(error));
  } finally {
    isBlocking.value = false;
  }
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isBlocking"
    :full-width="fullWidth"
    :title="title"
    @click.stop="requestBlock"
  >
    {{ buttonLabel }}
  </UiButton>

  <BlockLinkDialog
    v-model="isDialogOpen"
    :link="link"
    :loading="isBlocking"
    @confirm="confirmBlock"
  />
</template>
