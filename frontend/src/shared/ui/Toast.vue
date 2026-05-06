<script setup lang="ts">
export type UiToastItem = {
  id: string;
  type: "success" | "error" | "info" | "warning";
  title?: string;
  message: string;
  duration?: number;
};

withDefaults(
  defineProps<{
    items: UiToastItem[];
  }>(),
  {
    items: () => [],
  },
);

const emit = defineEmits<{
  remove: [id: string];
}>();

const onRemove = (id: string) => emit("remove", id);
</script>

<template>
  <Teleport to="body">
    <div class="ui-toast">
      <article v-for="item in items" :key="item.id" class="ui-toast__item" :class="`ui-toast__item--${item.type}`">
        <div class="ui-toast__content">
          <strong v-if="item.title" class="ui-toast__title">{{ item.title }}</strong>
          <p class="ui-toast__message">{{ item.message }}</p>
        </div>
        <button type="button" class="ui-toast__close" aria-label="Close toast" @click="onRemove(item.id)">x</button>
      </article>
    </div>
  </Teleport>
</template>

<style scoped>
.ui-toast {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 2100;
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: min(360px, calc(100vw - 32px));
}

.ui-toast__item {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  border-radius: var(--tl-radius-md);
  padding: 12px 14px;
  background: var(--tl-color-white);
  border: 1px solid #e6e1f0;
  box-shadow: 0 10px 20px rgb(20 16 38 / 12%);
}

.ui-toast__item--success {
  border-color: #73c090;
}

.ui-toast__item--error {
  border-color: #e58b95;
}

.ui-toast__item--info {
  border-color: #8ea8f6;
}

.ui-toast__item--warning {
  border-color: #eec783;
}

.ui-toast__content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ui-toast__title {
  color: var(--tl-color-text);
  font-size: 14px;
}

.ui-toast__message {
  margin: 0;
  color: var(--tl-color-text-muted);
  font-size: 13px;
}

.ui-toast__close {
  border: 0;
  background: transparent;
  color: var(--tl-color-text-muted);
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
}
</style>
