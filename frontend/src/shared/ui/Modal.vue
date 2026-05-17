<script setup lang="ts">
import { computed, onBeforeUnmount, watch } from "vue";
import { useI18n } from "@/shared/composables/useI18n";

type ModalWidth = "sm" | "md" | "lg";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    title?: string;
    description?: string;
    closeOnOverlay?: boolean;
    closeOnEsc?: boolean;
    width?: ModalWidth;
  }>(),
  {
    title: "",
    description: "",
    closeOnOverlay: true,
    closeOnEsc: true,
    width: "md",
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  close: [];
}>();

const isOpen = computed(() => props.modelValue);
const { t } = useI18n();

const closeModal = () => {
  emit("update:modelValue", false);
  emit("close");
};

const onOverlayClick = () => {
  if (props.closeOnOverlay) {
    closeModal();
  }
};

const onKeydown = (event: KeyboardEvent) => {
  if (!props.closeOnEsc || event.key !== "Escape" || !props.modelValue) {
    return;
  }
  closeModal();
};

watch(
  () => props.modelValue,
  (opened) => {
    if (opened) {
      document.addEventListener("keydown", onKeydown);
      return;
    }
    document.removeEventListener("keydown", onKeydown);
  },
  { immediate: true },
);

onBeforeUnmount(() => {
  document.removeEventListener("keydown", onKeydown);
});
</script>

<template>
  <Teleport to="body">
    <div v-if="isOpen" class="ui-modal">
      <button type="button" class="ui-modal__overlay" :aria-label="t('common.closeModal')" @click="onOverlayClick" />
      <section class="ui-modal__dialog" :class="`ui-modal__dialog--${width}`" role="dialog" aria-modal="true">
        <header v-if="title || description" class="ui-modal__header">
          <h3 v-if="title" class="ui-modal__title">{{ title }}</h3>
          <p v-if="description" class="ui-modal__description">{{ description }}</p>
        </header>

        <div class="ui-modal__body">
          <slot />
        </div>

        <footer v-if="$slots.footer" class="ui-modal__footer">
          <slot name="footer" />
        </footer>
      </section>
    </div>
  </Teleport>
</template>

<style scoped>
.ui-modal {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.ui-modal__overlay {
  position: absolute;
  inset: 0;
  border: 0;
  background: rgb(23 19 39 / 42%);
  cursor: pointer;
}

.ui-modal__dialog {
  position: relative;
  z-index: 1;
  width: 100%;
  border-radius: var(--tl-radius-lg);
  background: var(--tl-color-white);
  padding: 20px;
  box-shadow: 0 20px 35px rgb(20 16 38 / 20%);
}

.ui-modal__dialog--sm {
  max-width: 380px;
}

.ui-modal__dialog--md {
  max-width: 520px;
}

.ui-modal__dialog--lg {
  max-width: 720px;
}

.ui-modal__header {
  margin-bottom: 12px;
}

.ui-modal__title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: var(--tl-color-text);
}

.ui-modal__description {
  margin: 8px 0 0;
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.ui-modal__body {
  color: var(--tl-color-text);
}

.ui-modal__footer {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
