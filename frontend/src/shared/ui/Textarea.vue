<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    modelValue: string;
    label?: string;
    placeholder?: string;
    hint?: string;
    error?: string;
    disabled?: boolean;
    loading?: boolean;
    rows?: number;
    maxlength?: number;
  }>(),
  {
    label: "",
    placeholder: "",
    hint: "",
    error: "",
    disabled: false,
    loading: false,
    rows: 4,
    maxlength: undefined,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: string];
  blur: [event: FocusEvent];
  focus: [event: FocusEvent];
}>();

const isDisabled = computed(() => props.disabled || props.loading);
const hasError = computed(() => Boolean(props.error));

const model = computed({
  get: () => props.modelValue,
  set: (value: string) => emit("update:modelValue", value),
});

const onBlur = (event: FocusEvent) => emit("blur", event);
const onFocus = (event: FocusEvent) => emit("focus", event);
</script>

<template>
  <label class="ui-textarea" :class="{ 'is-disabled': isDisabled }">
    <span v-if="label" class="ui-textarea__label">{{ label }}</span>
    <textarea
      v-model="model"
      class="ui-textarea__field"
      :class="{ 'has-error': hasError, 'is-loading': loading }"
      :rows="rows"
      :maxlength="maxlength"
      :placeholder="placeholder"
      :disabled="isDisabled"
      @blur="onBlur"
      @focus="onFocus"
    />
    <span v-if="hasError" class="ui-textarea__error">{{ error }}</span>
    <span v-else-if="hint" class="ui-textarea__hint">{{ hint }}</span>
  </label>
</template>

<style scoped>
.ui-textarea {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.ui-textarea__label {
  font-size: 13px;
  font-weight: 700;
  color: var(--tl-color-text);
}

.ui-textarea__field {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd7e8;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
  color: var(--tl-color-text);
  font-size: 14px;
  font-family: var(--tl-font-family);
  resize: vertical;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, opacity 0.2s ease;
}

.ui-textarea__field:focus {
  outline: 0;
  border-color: var(--tl-color-primary);
  box-shadow: 0 0 0 2px rgb(109 74 255 / 18%);
}

.ui-textarea__field.has-error {
  border-color: var(--tl-color-danger);
}

.ui-textarea__field.is-loading {
  opacity: 0.75;
}

.ui-textarea__field::placeholder {
  color: var(--tl-color-text-muted);
}

.ui-textarea__hint {
  color: var(--tl-color-text-muted);
  font-size: 12px;
}

.ui-textarea__error {
  color: var(--tl-color-danger);
  font-size: 12px;
}

.ui-textarea.is-disabled .ui-textarea__field {
  opacity: 0.65;
}
</style>
