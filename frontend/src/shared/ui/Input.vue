<script setup lang="ts">
import { computed } from "vue";

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

const model = computed({
  get: () => props.modelValue,
  set: (value: string) => emit("update:modelValue", value),
});

const onBlur = (event: FocusEvent) => emit("blur", event);
const onFocus = (event: FocusEvent) => emit("focus", event);
</script>

<template>
  <label class="ui-input" :class="{ 'is-disabled': isDisabled }">
    <span v-if="label" class="ui-input__label">
      {{ label }}
      <span v-if="required" class="ui-input__required">*</span>
    </span>

    <span class="ui-input__control" :class="{ 'has-error': hasError, 'is-loading': loading }">
      <span v-if="$slots.prefix" class="ui-input__affix">
        <slot name="prefix" />
      </span>

      <input
        v-model="model"
        class="ui-input__field"
        :type="type"
        :placeholder="placeholder"
        :disabled="isDisabled"
        :required="required"
        :autocomplete="autocomplete"
        @blur="onBlur"
        @focus="onFocus"
      />

      <span v-if="errorInside && hasError" class="ui-input__error-inside">{{ error }}</span>
      <span v-else-if="$slots.suffix" class="ui-input__affix">
        <slot name="suffix" />
      </span>
    </span>

    <span v-if="hasError && !errorInside" class="ui-input__error">{{ error }}</span>
    <span v-else-if="hint" class="ui-input__hint">{{ hint }}</span>
  </label>
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

.ui-input__affix {
  display: inline-flex;
  align-items: center;
  color: var(--tl-color-text-muted);
  font-size: 13px;
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
