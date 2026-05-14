<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import {
  getLinkAnalytics,
  listRecentClicks,
  type ClickEvent,
  type LinkAnalyticsResponse,
} from "@/api/analytics";
import type { ApiClientError } from "@/api/types";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiButton, UiPageHeader, UiPageState, UiStatCard, UiTable, type UiTableColumn } from "@/shared/ui";

const route = useRoute();

const linkId = computed(() => {
  const rawId = route.params.id;
  return Array.isArray(rawId) ? rawId[0] : rawId;
});

const analytics = ref<LinkAnalyticsResponse | null>(null);
const recentClicks = ref<ClickEvent[]>([]);
const isLoading = ref(false);
const errorMessage = ref("");

const clicksColumns: UiTableColumn[] = [
  { key: "clickedAt", label: "Время", width: "24%" },
  { key: "referrer", label: "Источник", width: "30%" },
  { key: "userAgent", label: "User Agent", width: "46%" },
];

const numberFormatter = new Intl.NumberFormat("ru-RU");
const dateTimeFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const totalClicks = computed(() => numberFormatter.format(analytics.value?.totalClicks ?? 0));
const clicksLast24h = computed(() => numberFormatter.format(analytics.value?.clicksLast24h ?? 0));
const maxSeriesClicks = computed(() => Math.max(...(analytics.value?.series.map((point) => point.clicks) ?? [0]), 1));
const hasSeries = computed(() => Boolean(analytics.value?.series.length));

const lastClickedAt = computed(() => {
  if (!analytics.value?.lastClickedAt) {
    return "Переходов еще не было";
  }

  return formatDateTime(analytics.value.lastClickedAt);
});

const formatDateTime = (value: string) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "—";
  }

  return dateTimeFormatter.format(date);
};

const formatSource = (value?: string | null) => {
  const source = value?.trim();
  return source || "Прямой переход";
};

const formatUserAgent = (value?: string | null) => {
  const userAgent = value?.trim();
  return userAgent || "Не определен";
};

const getBarWidth = (clicks: number) => `${Math.max(6, Math.round((clicks / maxSeriesClicks.value) * 100))}%`;

const getErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401) {
      return "Сессия недействительна. Войдите заново, чтобы открыть аналитику ссылки.";
    }

    if (error.status === 404) {
      return "Ссылка не найдена или недоступна для просмотра.";
    }
  }

  return "Не удалось загрузить аналитику ссылки. Проверьте соединение и повторите попытку.";
};

const loadAnalytics = async () => {
  if (!linkId.value || isLoading.value) {
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";

  try {
    const [analyticsResponse, clicksResponse] = await Promise.all([
      getLinkAnalytics(linkId.value, { groupBy: "day" }),
      listRecentClicks(linkId.value, { limit: 20 }),
    ]);

    analytics.value = analyticsResponse;
    recentClicks.value = clicksResponse.items;
  } catch (error: unknown) {
    analytics.value = null;
    recentClicks.value = [];
    errorMessage.value = getErrorMessage(error);
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  void loadAnalytics();
});
</script>

<template>
  <section class="link-analytics-page">
    <UiPageHeader
      title="Аналитика ссылки"
      :subtitle="linkId ? `Подробная статистика переходов по ссылке ${linkId}.` : 'Ссылка не выбрана.'"
      :back-to="ROUTES.dashboard"
      back-label="Dashboard"
    />

    <UiPageState
      v-if="isLoading && !analytics"
      type="loading"
      title="Загружаем аналитику"
      description="Получаем метрики, динамику трафика и последние переходы."
    />

    <UiPageState
      v-else-if="errorMessage"
      type="error"
      title="Аналитика недоступна"
      :description="errorMessage"
      action-text="Повторить"
      @action="loadAnalytics"
    />

    <UiPageState
      v-else-if="!linkId"
      type="not-found"
      title="Ссылка не выбрана"
      description="Откройте аналитику из списка ссылок на dashboard."
      :action-to="ROUTES.dashboard"
      action-text="Вернуться на dashboard"
    />

    <div v-else-if="analytics" class="link-analytics-page__content">
      <section class="link-analytics-page__summary" aria-label="Основные метрики">
        <UiStatCard title="Всего переходов" :value="totalClicks" />
        <UiStatCard title="Переходы за 24 часа" :value="clicksLast24h" />
        <UiStatCard title="Последний переход" :value="lastClickedAt" />
      </section>

      <section class="link-analytics-page__traffic" aria-labelledby="link-analytics-traffic-title">
        <div class="link-analytics-page__section-header">
          <h2 id="link-analytics-traffic-title" class="link-analytics-page__section-title">Трафик</h2>
          <UiButton variant="secondary" size="sm" :loading="isLoading" @click="loadAnalytics">Обновить</UiButton>
        </div>

        <div v-if="hasSeries" class="link-analytics-page__chart">
          <div v-for="point in analytics.series" :key="point.periodStart" class="link-analytics-page__bar-row">
            <span class="link-analytics-page__bar-label">{{ formatDateTime(point.periodStart) }}</span>
            <div class="link-analytics-page__bar-track" aria-hidden="true">
              <span class="link-analytics-page__bar" :style="{ width: getBarWidth(point.clicks) }" />
            </div>
            <strong class="link-analytics-page__bar-value">{{ numberFormatter.format(point.clicks) }}</strong>
          </div>
        </div>

        <UiPageState
          v-else
          type="empty"
          title="Трафика пока нет"
          description="Когда по ссылке появятся переходы, здесь отобразится динамика по дням."
        />
      </section>

      <section class="link-analytics-page__clicks" aria-labelledby="link-analytics-clicks-title">
        <div class="link-analytics-page__section-header">
          <h2 id="link-analytics-clicks-title" class="link-analytics-page__section-title">Последние переходы</h2>
        </div>

        <UiTable :columns="clicksColumns" :rows="recentClicks" empty-text="Переходов пока нет.">
          <template #empty>
            <UiPageState
              type="empty"
              title="Переходов пока нет"
              description="Последние переходы появятся после первых открытий короткой ссылки."
            />
          </template>

          <template #cell="{ row, column }">
            <span v-if="column.key === 'clickedAt'">{{ formatDateTime(row.clickedAt) }}</span>
            <a
              v-else-if="column.key === 'referrer' && row.referrer"
              class="link-analytics-page__referrer"
              :href="row.referrer"
              target="_blank"
              rel="noreferrer"
              :title="row.referrer"
            >
              {{ formatSource(row.referrer) }}
            </a>
            <span v-else-if="column.key === 'referrer'">{{ formatSource(row.referrer) }}</span>
            <span v-else-if="column.key === 'userAgent'" class="link-analytics-page__user-agent">
              {{ formatUserAgent(row.userAgent) }}
            </span>
            <span v-else>{{ row[column.key] }}</span>
          </template>
        </UiTable>
      </section>
    </div>
  </section>
</template>

<style scoped>
.link-analytics-page {
  width: 100%;
}

.link-analytics-page__content {
  display: flex;
  flex-direction: column;
  gap: 26px;
}

.link-analytics-page__summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
}

.link-analytics-page__traffic,
.link-analytics-page__clicks {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-top: 22px;
  border-top: 1px solid rgb(37 31 63 / 10%);
}

.link-analytics-page__section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.link-analytics-page__section-title {
  margin: 0;
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.link-analytics-page__chart {
  display: grid;
  gap: 10px;
}

.link-analytics-page__bar-row {
  display: grid;
  grid-template-columns: minmax(150px, 190px) minmax(160px, 1fr) 56px;
  align-items: center;
  gap: 12px;
  color: var(--tl-color-text);
  font-size: 14px;
}

.link-analytics-page__bar-label {
  color: var(--tl-color-text-muted);
}

.link-analytics-page__bar-track {
  height: 14px;
  overflow: hidden;
  border-radius: 999px;
  background: var(--tl-color-surface-muted);
}

.link-analytics-page__bar {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: var(--tl-color-primary);
}

.link-analytics-page__bar-value {
  text-align: right;
}

.link-analytics-page__referrer,
.link-analytics-page__user-agent {
  display: inline-block;
  max-width: 360px;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: bottom;
  white-space: nowrap;
}

.link-analytics-page__referrer {
  color: var(--tl-color-primary);
}

@media (max-width: 767px) {
  .link-analytics-page__summary {
    grid-template-columns: 1fr;
  }

  .link-analytics-page__section-header,
  .link-analytics-page__bar-row {
    align-items: stretch;
    grid-template-columns: 1fr;
  }

  .link-analytics-page__bar-value {
    text-align: left;
  }

  .link-analytics-page__referrer,
  .link-analytics-page__user-agent {
    max-width: 220px;
  }
}
</style>
