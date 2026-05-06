<script setup lang="ts">
import { computed } from "vue";

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";
type ButtonSize = "sm" | "md" | "lg";
type NativeButtonType = "button" | "submit" | "reset";

const props = withDefaults(
  defineProps<{
    variant?: ButtonVariant;
    size?: ButtonSize;
    type?: NativeButtonType;
    disabled?: boolean;
    loading?: boolean;
    fullWidth?: boolean;
  }>(),
  {
    variant: "primary",
    size: "md",
    type: "button",
    disabled: false,
    loading: false,
    fullWidth: false,
  },
);

const isDisabled = computed(() => props.disabled || props.loading);
</script>

<template>
  <button
    class="ui-button"
    :class="[
      `ui-button--${variant}`,
      `ui-button--${size}`,
      { 'is-loading': loading, 'is-full-width': fullWidth },
    ]"
    :type="type"
    :disabled="isDisabled"
    :aria-busy="loading"
  >
    <span v-if="loading" class="ui-button__spinner" aria-hidden="true" />
    <span v-else-if="$slots.iconLeft" class="ui-button__icon">
      <slot name="iconLeft" />
    </span>

    <span class="ui-button__content">
      <slot />
    </span>

    <span v-if="!loading && $slots.iconRight" class="ui-button__icon">
      <slot name="iconRight" />
    </span>
  </button>
</template>

<style scoped>
.ui-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: fit-content;
  border: 1px solid transparent;
  border-radius: var(--tl-radius-lg);
  font-family: var(--tl-font-family);
  font-weight: 600;
  line-height: 1;
  cursor: pointer;
  transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease, opacity 0.2s ease;
}

.ui-button.is-full-width {
  width: 100%;
}

.ui-button:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.ui-button--sm {
  min-height: 34px;
  padding: 8px 12px;
  font-size: 13px;
}

.ui-button--md {
  min-height: 40px;
  padding: 10px 16px;
  font-size: 14px;
}

.ui-button--lg {
  min-height: 46px;
  padding: 12px 18px;
  font-size: 15px;
}

.ui-button--primary {
  background: var(--tl-color-primary);
  color: var(--tl-color-white);
}

.ui-button--primary:hover:not(:disabled) {
  background: var(--tl-color-primary-hover);
}

.ui-button--secondary {
  background: var(--tl-color-surface-muted);
  border-color: #e4dfef;
  color: var(--tl-color-text);
}

.ui-button--secondary:hover:not(:disabled) {
  background: #e8e2f4;
}

.ui-button--danger {
  background: var(--tl-color-danger);
  color: var(--tl-color-white);
}

.ui-button--danger:hover:not(:disabled) {
  background: var(--tl-color-danger-hover);
}

.ui-button--ghost {
  background: transparent;
  color: var(--tl-color-primary);
}

.ui-button--ghost:hover:not(:disabled) {
  background: rgb(109 74 255 / 10%);
}

.ui-button__icon {
  display: inline-flex;
  align-items: center;
}

.ui-button__spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgb(255 255 255 / 45%);
  border-top-color: currentcolor;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.ui-button--secondary .ui-button__spinner,
.ui-button--ghost .ui-button__spinner {
  border-color: rgb(37 31 63 / 25%);
  border-top-color: currentcolor;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
