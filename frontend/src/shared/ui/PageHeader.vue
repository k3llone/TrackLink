<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";

const props = withDefaults(
  defineProps<{
    title?: string;
    subtitle?: string;
    backTo?: string;
    backLabel?: string;
  }>(),
  {
    title: "",
    subtitle: "",
    backTo: "",
    backLabel: "",
  },
);

const { t } = useI18n();
const resolvedBackLabel = computed(() => props.backLabel || t("common.back"));
</script>

<template>
  <header class="ui-page-header">
    <div class="ui-page-header__main">
      <a v-if="backTo" :href="backTo" class="ui-page-header__back">&larr; {{ resolvedBackLabel }}</a>
      <h1 v-if="title" class="ui-page-header__title">{{ title }}</h1>
      <p v-if="subtitle" class="ui-page-header__subtitle">{{ subtitle }}</p>
      <slot />
    </div>

    <div v-if="$slots.actions" class="ui-page-header__actions">
      <slot name="actions" />
    </div>
  </header>
</template>

<style scoped>
.ui-page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
}

.ui-page-header__main {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ui-page-header__back {
  width: fit-content;
  color: var(--tl-color-primary);
  text-decoration: none;
  font-weight: 600;
  font-size: 13px;
}

.ui-page-header__title {
  margin: 0;
  color: var(--tl-color-primary);
  font-size: 28px;
  line-height: 1.2;
}

.ui-page-header__subtitle {
  margin: 0;
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.ui-page-header__actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
</style>
