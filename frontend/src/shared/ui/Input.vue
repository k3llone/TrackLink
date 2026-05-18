<script setup lang="ts">
import { computed, ref, useId } from "vue";
import { useI18n } from "@/shared/composables/useI18n";

type InputType = "text" | "email" | "password" | "url" | "search";

const props = withDefaults(
  defineProps<{
    modelValue: string;
    label?: string;
    type?: InputType;
    placeholder?: string;
    hint?: string;
    error?: string;
    disabled?: boolean;
    readonly?: boolean;
    loading?: boolean;
    required?: boolean;
    autocomplete?: string;
    errorInside?: boolean;
  }>(),
  {
    label: "",
    type: "text",
    placeholder: "",
    hint: "",
    error: "",
    disabled: false,
    readonly: false,
    loading: false,
    required: false,
    autocomplete: "off",
    errorInside: false,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: string];
  blur: [event: FocusEvent];
  focus: [event: FocusEvent];
}>();

const hasError = computed(() => Boolean(props.error));
const isDisabled = computed(() => props.disabled || props.loading);
const isPasswordInput = computed(() => props.type === "password");
const isPasswordVisible = ref(false);
const inputId = useId();
const { t } = useI18n();

const model = computed({
  get: () => props.modelValue,
  set: (value: string) => emit("update:modelValue", value),
});

const inputType = computed<InputType>(() => {
  if (!isPasswordInput.value) {
    return props.type;
  }

  return isPasswordVisible.value ? "text" : "password";
});
const passwordToggleLabel = computed(() =>
  isPasswordVisible.value ? t("common.hidePassword") : t("common.showPassword"),
);

const onBlur = (event: FocusEvent) => emit("blur", event);
const onFocus = (event: FocusEvent) => emit("focus", event);
const togglePasswordVisibility = () => {
  if (isDisabled.value || props.readonly) {
    return;
  }

  isPasswordVisible.value = !isPasswordVisible.value;
};
</script>

<template>
  <div class="ui-input" :class="{ 'is-disabled': isDisabled }">
    <label v-if="label" class="ui-input__label" :for="inputId">
      {{ label }}
      <span v-if="required" class="ui-input__required">*</span>
    </label>

    <span class="ui-input__control" :class="{ 'has-error': hasError, 'is-loading': loading }">
      <span v-if="$slots.prefix" class="ui-input__affix">
        <slot name="prefix" />
      </span>

      <input
        :id="inputId"
        v-model="model"
        class="ui-input__field"
        :type="inputType"
        :placeholder="placeholder"
        :disabled="isDisabled"
        :readonly="readonly"
        :required="required"
        :autocomplete="autocomplete"
        @blur="onBlur"
        @focus="onFocus"
      />

      <span v-if="errorInside && hasError" class="ui-input__error-inside">{{ error }}</span>
      <button
        v-if="isPasswordInput"
        class="ui-input__password-toggle"
        type="button"
        :aria-label="passwordToggleLabel"
        :aria-pressed="isPasswordVisible"
        :disabled="isDisabled || readonly"
        :title="passwordToggleLabel"
        @mousedown.prevent
        @click="togglePasswordVisibility"
      >
        <svg v-if="isPasswordVisible" viewBox="0 0 24 24" focusable="false" aria-hidden="true">
          <path
            d="M2.2 3.6 3.6 2.2l18.2 18.2-1.4 1.4-3.2-3.2A11.3 11.3 0 0 1 12 20C7 20 2.7 16.9 1 12c.8-2.2 2.1-4.1 3.8-5.5L2.2 3.6Zm7.1 7.1A3 3 0 0 0 12 15a3 3 0 0 0 1.3-.3l-4-4Zm2.4-3.6 5.9 5.9.1-1A5.7 5.7 0 0 0 12 6.3l-.3.8ZM12 4c5 0 9.3 3.1 11 8a13 13 0 0 1-2.7 4.4l-2-2A9 9 0 0 0 20.8 12C19.3 8.6 15.9 6 12 6c-1 0-1.9.2-2.8.5L7.6 4.9A11.6 11.6 0 0 1 12 4Z"
            fill="currentColor"
          />
        </svg>
        <svg v-else viewBox="0 0 24 24" focusable="false" aria-hidden="true">
          <path
            d="M12 4c5 0 9.3 3.1 11 8-1.7 4.9-6 8-11 8S2.7 16.9 1 12c1.7-4.9 6-8 11-8Zm0 2c-3.9 0-7.3 2.3-8.8 6 1.5 3.7 4.9 6 8.8 6s7.3-2.3 8.8-6C19.3 8.3 15.9 6 12 6Zm0 2.5a3.5 3.5 0 1 1 0 7 3.5 3.5 0 0 1 0-7Zm0 2a1.5 1.5 0 1 0 0 3 1.5 1.5 0 0 0 0-3Z"
            fill="currentColor"
          />
        </svg>
      </button>
      <span v-else-if="$slots.suffix" class="ui-input__affix">
        <slot name="suffix" />
      </span>
    </span>

    <span v-if="hasError && !errorInside" class="ui-input__error">{{ error }}</span>
    <span v-else-if="hint" class="ui-input__hint">{{ hint }}</span>
  </div>
</template>

<style scoped>
.ui-input {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
  color: var(--tl-color-text);
}

.ui-input__label {
  font-size: 13px;
  font-weight: 700;
}

.ui-input__required {
  margin-left: 4px;
  color: var(--tl-color-danger);
}

.ui-input__control {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 44px;
  padding: 0 12px;
  border: 1px solid #ddd7e8;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
  transition: border-color 0.2s ease, box-shadow 0.2s ease, opacity 0.2s ease;
}

.ui-input__control:focus-within {
  border-color: var(--tl-color-primary);
  box-shadow: 0 0 0 2px rgb(109 74 255 / 18%);
}

.ui-input__control.has-error {
  border-color: var(--tl-color-danger);
}

.ui-input__control.is-loading {
  opacity: 0.75;
}

.ui-input__field {
  width: 100%;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--tl-color-text);
  font-size: 14px;
  font-family: var(--tl-font-family);
}

.ui-input__field::placeholder {
  color: var(--tl-color-text-muted);
}

.ui-input__field::-ms-reveal,
.ui-input__field::-ms-clear {
  display: none;
}

.ui-input__affix {
  display: inline-flex;
  align-items: center;
  color: var(--tl-color-text-muted);
  font-size: 13px;
}

.ui-input__password-toggle {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  margin-right: -4px;
  border: 0;
  border-radius: var(--tl-radius-sm);
  background: transparent;
  color: var(--tl-color-text-muted);
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease, opacity 0.2s ease;
}

.ui-input__password-toggle:hover:not(:disabled) {
  background: rgb(109 74 255 / 10%);
  color: var(--tl-color-primary);
}

.ui-input__password-toggle:focus-visible {
  outline: 2px solid rgb(109 74 255 / 35%);
  outline-offset: 2px;
}

.ui-input__password-toggle:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.ui-input__password-toggle svg {
  width: 18px;
  height: 18px;
}

.ui-input__error-inside {
  color: var(--tl-color-danger);
  font-size: 12px;
  white-space: nowrap;
}

.ui-input__hint {
  color: var(--tl-color-text-muted);
  font-size: 12px;
}

.ui-input__error {
  color: var(--tl-color-danger);
  font-size: 12px;
}

.ui-input.is-disabled .ui-input__control {
  opacity: 0.65;
}
</style>
