<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    disabled?: boolean;
    loading?: boolean;
    label?: string;
    description?: string;
  }>(),
  {
    disabled: false,
    loading: false,
    label: "",
    description: "",
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  change: [value: boolean];
}>();

const isDisabled = computed(() => props.disabled || props.loading);

const onToggle = () => {
  if (isDisabled.value) {
    return;
  }

  const nextValue = !props.modelValue;
  emit("update:modelValue", nextValue);
  emit("change", nextValue);
};
</script>

<template>
  <button
    type="button"
    class="ui-switch"
    :class="{ 'is-checked': modelValue, 'is-disabled': isDisabled }"
    :disabled="isDisabled"
    @click="onToggle"
  >
    <span class="ui-switch__content">
      <span v-if="label" class="ui-switch__label">{{ label }}</span>
      <span v-if="description" class="ui-switch__description">{{ description }}</span>
    </span>
    <span class="ui-switch__track" aria-hidden="true">
      <span class="ui-switch__thumb" />
    </span>
  </button>
</template>

<style scoped>
.ui-switch {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  width: 100%;
  padding: 8px 0;
  border: 0;
  background: transparent;
  font-family: var(--tl-font-family);
  text-align: left;
  cursor: pointer;
}

.ui-switch.is-disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.ui-switch__content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.ui-switch__label {
  font-size: 14px;
  font-weight: 700;
  color: var(--tl-color-text);
}

.ui-switch__description {
  font-size: 12px;
  color: var(--tl-color-text-muted);
}

.ui-switch__track {
  position: relative;
  display: inline-flex;
  align-items: center;
  width: 44px;
  height: 24px;
  padding: 2px;
  border-radius: 999px;
  background: #d9d2ea;
  transition: background-color 0.2s ease;
}

.ui-switch.is-checked .ui-switch__track {
  background: var(--tl-color-primary);
}

.ui-switch__thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--tl-color-white);
  box-shadow: 0 1px 3px rgb(0 0 0 / 20%);
  transform: translateX(0);
  transition: transform 0.2s ease;
}

.ui-switch.is-checked .ui-switch__thumb {
  transform: translateX(20px);
}
</style>
