<script setup lang="ts">
import { computed, ref } from "vue";
import { deleteLink } from "@/api/links";
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
    label: "Удалить",
  },
);

const emit = defineEmits<{
  deleted: [linkId: string];
}>();

const toast = useToast();
const isConfirmOpen = ref(false);
const isDeleting = ref(false);

const isDisabled = computed(() => props.disabled || isDeleting.value || !props.linkId);
const confirmDescription = computed(() => {
  const shortUrl = props.shortUrl.trim();

  if (shortUrl) {
    return `Ссылка ${shortUrl} будет удалена из обычного списка. Это действие нельзя отменить.`;
  }

  return "Ссылка будет удалена из обычного списка. Это действие нельзя отменить.";
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
    toast.success("Ссылка удалена.");
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
    {{ label }}
  </UiButton>

  <UiConfirmDialog
    v-model="isConfirmOpen"
    title="Удалить ссылку?"
    :description="confirmDescription"
    confirm-text="Удалить"
    cancel-text="Отмена"
    :loading="isDeleting"
    @confirm="confirmDelete"
  />
</template>
