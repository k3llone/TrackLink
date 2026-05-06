<script setup lang="ts">
import UiButton from "./Button.vue";
import UiModal from "./Modal.vue";

type ConfirmVariant = "primary" | "danger";

withDefaults(
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
    confirmText: "Confirm",
    cancelText: "Cancel",
    variant: "danger",
    loading: false,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  confirm: [];
  cancel: [];
}>();

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
      <UiButton variant="secondary" :disabled="loading" @click="onCancel">{{ cancelText }}</UiButton>
      <UiButton :variant="variant" :loading="loading" @click="onConfirm">{{ confirmText }}</UiButton>
    </template>
  </UiModal>
</template>
