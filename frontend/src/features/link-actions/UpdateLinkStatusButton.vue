<script setup lang="ts">
import { computed, ref } from "vue";
import { updateLinkStatus } from "@/api/links";
import type { Link, UpdateLinkStatus } from "@/entities/link/link.types";
import { useI18n } from "@/shared/composables/useI18n";
import { useToast } from "@/shared/composables/useToast";
import { UiButton } from "@/shared/ui";
import { getUpdateStatusErrorMessage } from "./linkActionErrors";

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

const props = withDefaults(
  defineProps<{
    link: Link;
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
  }>(),
  {
    variant: "ghost",
    size: "sm",
    disabled: false,
  },
);

const emit = defineEmits<{
  updated: [link: Link];
}>();

const toast = useToast();
const { t } = useI18n();
const isUpdating = ref(false);

const nextStatus = computed<UpdateLinkStatus | null>(() => {
  if (props.link.status === "active") {
    return "inactive";
  }

  if (props.link.status === "inactive") {
    return "active";
  }

  return null;
});

const canUpdateStatus = computed(() => Boolean(nextStatus.value));
const isDisabled = computed(() => props.disabled || isUpdating.value || !canUpdateStatus.value);
const actionLabel = computed(() => {
  if (nextStatus.value === "active") {
    return t("common.activate");
  }

  if (nextStatus.value === "inactive") {
    return t("common.deactivate");
  }

  return t("linkActions.status.unavailable");
});
const title = computed(() => (canUpdateStatus.value ? actionLabel.value : t("linkActions.status.blockedOrDeletedTitle")));

const getSuccessMessage = (status: UpdateLinkStatus) =>
  status === "active" ? t("linkActions.status.activateSuccess") : t("linkActions.status.deactivateSuccess");

const onUpdateStatus = async () => {
  if (!nextStatus.value || isUpdating.value) {
    return;
  }

  isUpdating.value = true;

  try {
    const updatedLink = await updateLinkStatus(props.link.id, { status: nextStatus.value });
    toast.success(getSuccessMessage(updatedLink.status as UpdateLinkStatus));
    emit("updated", updatedLink);
  } catch (error: unknown) {
    toast.error(getUpdateStatusErrorMessage(error));
  } finally {
    isUpdating.value = false;
  }
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isUpdating"
    :title="title"
    @click.stop="onUpdateStatus"
  >
    {{ actionLabel }}
  </UiButton>
</template>
