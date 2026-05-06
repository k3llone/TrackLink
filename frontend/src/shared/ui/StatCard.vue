<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string;
    value: string | number;
    hint?: string;
    loading?: boolean;
  }>(),
  {
    hint: "",
    loading: false,
  },
);
</script>

<template>
  <article class="ui-stat-card">
    <header class="ui-stat-card__header">
      <span class="ui-stat-card__title">{{ title }}</span>
      <span v-if="$slots.icon" class="ui-stat-card__icon">
        <slot name="icon" />
      </span>
    </header>

    <div v-if="loading" class="ui-stat-card__skeleton" />
    <p v-else class="ui-stat-card__value">{{ value }}</p>
    <p v-if="hint" class="ui-stat-card__hint">{{ hint }}</p>
  </article>
</template>

<style scoped>
.ui-stat-card {
  border-radius: var(--tl-radius-lg);
  background: var(--tl-color-white);
  border: 1px solid #ece7f4;
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ui-stat-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.ui-stat-card__title {
  font-size: 13px;
  color: var(--tl-color-text-muted);
  font-weight: 700;
}

.ui-stat-card__icon {
  display: inline-flex;
  color: var(--tl-color-primary);
}

.ui-stat-card__value {
  margin: 0;
  font-size: 30px;
  line-height: 1.1;
  color: var(--tl-color-text);
  font-weight: 700;
}

.ui-stat-card__hint {
  margin: 0;
  font-size: 12px;
  color: var(--tl-color-text-muted);
}

.ui-stat-card__skeleton {
  height: 34px;
  width: 55%;
  border-radius: 8px;
  background: linear-gradient(90deg, #ece7f4 20%, #f4f1fa 45%, #ece7f4 80%);
  background-size: 200% 100%;
  animation: pulse 1.4s ease-in-out infinite;
}

@keyframes pulse {
  0% {
    background-position: 100% 0;
  }
  100% {
    background-position: -100% 0;
  }
}
</style>
