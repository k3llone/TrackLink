<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import UiButton from "./Button.vue";

type PageStateType = "loading" | "empty" | "error" | "forbidden" | "not-found";

const props = withDefaults(
  defineProps<{
    type: PageStateType;
    title?: string;
    description?: string;
    actionText?: string;
    actionTo?: string;
  }>(),
  {
    title: "",
    description: "",
    actionText: "",
    actionTo: "",
  },
);

const emit = defineEmits<{
  action: [];
}>();

const { t } = useI18n();

const presets: Record<PageStateType, { title: () => string; description: () => string }> = {
  loading: {
    title: () => t("pageState.loading.title"),
    description: () => t("pageState.loading.description"),
  },
  empty: {
    title: () => t("pageState.empty.title"),
    description: () => t("pageState.empty.description"),
  },
  error: {
    title: () => t("pageState.error.title"),
    description: () => t("pageState.error.description"),
  },
  forbidden: {
    title: () => t("pageState.forbidden.title"),
    description: () => t("pageState.forbidden.description"),
  },
  "not-found": {
    title: () => t("pageState.notFound.title"),
    description: () => t("pageState.notFound.description"),
  },
};

const resolvedTitle = computed(() => props.title || presets[props.type].title());
const resolvedDescription = computed(() => props.description || presets[props.type].description());

const onAction = () => emit("action");
</script>

<template>
  <section class="ui-page-state" :class="`ui-page-state--${type}`">
    <div v-if="type === 'loading'" class="ui-page-state__spinner" aria-hidden="true" />
    <h2 class="ui-page-state__title">{{ resolvedTitle }}</h2>
    <p class="ui-page-state__description">{{ resolvedDescription }}</p>

    <a v-if="actionText && actionTo" class="ui-page-state__link" :href="actionTo">
      <UiButton variant="primary">{{ actionText }}</UiButton>
    </a>
    <UiButton v-else-if="actionText" variant="primary" @click="onAction">
      {{ actionText }}
    </UiButton>
  </section>
</template>

<style scoped>
.ui-page-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 10px;
  padding: 40px 16px;
}

.ui-page-state__title {
  margin: 0;
  color: var(--tl-color-text);
  font-size: 22px;
}

.ui-page-state__description {
  margin: 0;
  color: var(--tl-color-text-muted);
  font-size: 14px;
  max-width: 420px;
}

.ui-page-state__spinner {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 3px solid rgb(109 74 255 / 25%);
  border-top-color: var(--tl-color-primary);
  animation: spin 0.7s linear infinite;
}

.ui-page-state__link {
  text-decoration: none;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
