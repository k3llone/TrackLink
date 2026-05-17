<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import UiButton from "./Button.vue";
import UiModal from "./Modal.vue";

type ConfirmVariant = "primary" | "danger";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    title: string;
    description?: string;
    confirmText?: string;
    cancelText?: string;
    variant?: ConfirmVariant;
    loading?: boolean;
  }>(),
  {
    description: "",
    confirmText: "",
    cancelText: "",
    variant: "danger",
    loading: false,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  confirm: [];
  cancel: [];
}>();

const { t } = useI18n();
const resolvedConfirmText = computed(() => props.confirmText || t("common.confirm"));
const resolvedCancelText = computed(() => props.cancelText || t("common.cancel"));

const onClose = () => emit("update:modelValue", false);
const onCancel = () => {
  emit("cancel");
  emit("update:modelValue", false);
};
const onConfirm = () => emit("confirm");
</script>

<template>
  <UiModal :model-value="modelValue" :title="title" :description="description" @update:model-value="onClose">
    <template #footer>
      <UiButton variant="secondary" :disabled="loading" @click="onCancel">{{ resolvedCancelText }}</UiButton>
      <UiButton :variant="variant" :loading="loading" @click="onConfirm">{{ resolvedConfirmText }}</UiButton>
    </template>
  </UiModal>
</template>
