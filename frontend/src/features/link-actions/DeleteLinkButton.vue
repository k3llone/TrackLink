<script setup lang="ts">
import { computed, ref } from "vue";
import { deleteLink } from "@/api/links";
import { useI18n } from "@/shared/composables/useI18n";
import { useToast } from "@/shared/composables/useToast";
import { UiButton, UiConfirmDialog } from "@/shared/ui";
import { getDeleteLinkErrorMessage } from "./linkActionErrors";

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

const props = withDefaults(
  defineProps<{
    linkId: string;
    shortUrl?: string;
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
    label?: string;
  }>(),
  {
    shortUrl: "",
    variant: "ghost",
    size: "sm",
    disabled: false,
    label: "",
  },
);

const emit = defineEmits<{
  deleted: [linkId: string];
}>();

const toast = useToast();
const { t } = useI18n();
const isConfirmOpen = ref(false);
const isDeleting = ref(false);

const isDisabled = computed(() => props.disabled || isDeleting.value || !props.linkId);
const buttonLabel = computed(() => props.label || t("common.delete"));
const confirmDescription = computed(() => {
  const shortUrl = props.shortUrl.trim();

  if (shortUrl) {
    return t("linkActions.delete.descriptionWithUrl", { shortUrl });
  }

  return t("linkActions.delete.description");
});

const requestDelete = () => {
  if (isDisabled.value) {
    return;
  }

  isConfirmOpen.value = true;
};

const confirmDelete = async () => {
  if (isDisabled.value) {
    return;
  }

  isDeleting.value = true;

  try {
    await deleteLink(props.linkId);
    isConfirmOpen.value = false;
    toast.success(t("linkActions.delete.success"));
    emit("deleted", props.linkId);
  } catch (error: unknown) {
    toast.error(getDeleteLinkErrorMessage(error));
  } finally {
    isDeleting.value = false;
  }
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isDeleting"
    @click.stop="requestDelete"
  >
    {{ buttonLabel }}
  </UiButton>

  <UiConfirmDialog
    v-model="isConfirmOpen"
    :title="t('linkActions.delete.confirmTitle')"
    :description="confirmDescription"
    :confirm-text="t('common.delete')"
    :cancel-text="t('common.cancel')"
    :loading="isDeleting"
    @confirm="confirmDelete"
  />
</template>
