<script setup lang="ts">
import { computed } from "vue";
import type { LinkAnalyticsResponse } from "@/api/analytics";
import type { Link } from "@/entities/link/link.types";
import { UiButton, UiModal, UiPageState, UiStatCard } from "@/shared/ui";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    link: Link | null;
    analytics: LinkAnalyticsResponse | null;
    loading?: boolean;
    errorMessage?: string;
  }>(),
  {
    loading: false,
    errorMessage: "",
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  retry: [];
}>();

const numberFormatter = new Intl.NumberFormat("ru-RU");
const dateTimeFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});

const modalTitle = computed(() => (props.link ? `Аналитика: ${props.link.shortUrl || props.link.code}` : "Аналитика"));
const totalClicks = computed(() => numberFormatter.format(props.analytics?.totalClicks ?? 0));
const clicksLast24h = computed(() => numberFormatter.format(props.analytics?.clicksLast24h ?? 0));
const series = computed(() => props.analytics?.series.slice(0, 6) ?? []);
const hasSeries = computed(() => series.value.length > 0);
const lastClickedAt = computed(() => {
  if (!props.analytics?.lastClickedAt) {
    return "Переходов еще не было";
  }

  const date = new Date(props.analytics.lastClickedAt);

  if (Number.isNaN(date.getTime())) {
    return "Дата недоступна";
  }

  return dateTimeFormatter.format(date);
});

const formatPeriod = (value: string) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "—";
  }

  return dateTimeFormatter.format(date);
};

const closeModal = () => emit("update:modelValue", false);
const retry = () => emit("retry");
</script>

<template>
  <UiModal
    :model-value="modelValue"
    :title="modalTitle"
    description="Базовая статистика переходов по выбранной ссылке."
    width="lg"
    @update:model-value="closeModal"
  >
    <UiPageState
      v-if="loading"
      type="loading"
      title="Загружаем аналитику"
      description="Получаем статистику выбранной ссылки."
    />

    <UiPageState
      v-else-if="errorMessage"
      type="error"
      title="Аналитика недоступна"
      :description="errorMessage"
      action-text="Повторить"
      @action="retry"
    />

    <div v-else-if="analytics" class="link-analytics-modal__content">
      <div class="link-analytics-modal__stats">
        <UiStatCard title="Всего переходов" :value="totalClicks" />
        <UiStatCard title="За 24 часа" :value="clicksLast24h" />
        <UiStatCard title="Последний переход" :value="lastClickedAt" />
      </div>

      <section class="link-analytics-modal__series" aria-labelledby="link-analytics-series-title">
        <h4 id="link-analytics-series-title" class="link-analytics-modal__series-title">Динамика переходов</h4>

        <ul v-if="hasSeries" class="link-analytics-modal__series-list">
          <li v-for="point in series" :key="point.periodStart" class="link-analytics-modal__series-item">
            <span>{{ formatPeriod(point.periodStart) }}</span>
            <strong>{{ numberFormatter.format(point.clicks) }}</strong>
          </li>
        </ul>

        <p v-else class="link-analytics-modal__empty-series">По этой ссылке пока нет данных по периодам.</p>
      </section>
    </div>

    <template #footer>
      <UiButton variant="secondary" type="button" @click="closeModal">Закрыть</UiButton>
    </template>
  </UiModal>
</template>

<style scoped>
.link-analytics-modal__content {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.link-analytics-modal__stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.link-analytics-modal__series {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.link-analytics-modal__series-title {
  margin: 0;
  color: var(--tl-color-text);
  font-size: 16px;
}

.link-analytics-modal__series-list {
  display: grid;
  gap: 8px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.link-analytics-modal__series-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
  padding: 10px 12px;
  color: var(--tl-color-text);
  font-size: 14px;
}

.link-analytics-modal__empty-series {
  margin: 0;
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

@media (max-width: 767px) {
  .link-analytics-modal__stats {
    grid-template-columns: 1fr;
  }
}
</style>
