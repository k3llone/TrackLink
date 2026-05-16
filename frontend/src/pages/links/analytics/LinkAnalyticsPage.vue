<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  getLinkAnalytics,
  listRecentClicks,
  type ClickEvent,
  type LinkAnalyticsResponse,
} from "@/api/analytics";
import { findOwnLinkById } from "@/api/links";
import type { ApiClientError } from "@/api/types";
import type { Link } from "@/entities/link/link.types";
import CopyShortUrlButton from "@/features/link-actions/CopyShortUrlButton.vue";
import DeleteLinkButton from "@/features/link-actions/DeleteLinkButton.vue";
import UpdateLinkStatusButton from "@/features/link-actions/UpdateLinkStatusButton.vue";
import { ROUTES } from "@/shared/lib/routes/paths";
import {
  UiPageHeader,
  UiPageState,
  UiStatCard,
  UiTable,
  type UiTableColumn,
} from "@/shared/ui";

const route = useRoute();
const router = useRouter();

const linkId = computed(() => {
  const rawId = route.params.id;
  return Array.isArray(rawId) ? rawId[0] : rawId;
});

const analytics = ref<LinkAnalyticsResponse | null>(null);
const link = ref<Link | null>(null);
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
const pageSubtitle = computed(() => {
  if (link.value) {
    return `Подробная статистика переходов по ссылке ${link.value.shortUrl}.`;
  }

  return linkId.value ? `Подробная статистика переходов по ссылке ${linkId.value}.` : "Ссылка не выбрана.";
});

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

const getLinkNotFoundMessage = () => "Ссылка не найдена или недоступна для управления.";

const loadAnalytics = async () => {
  if (!linkId.value || isLoading.value) {
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";

  try {
    const [linkResponse, analyticsResponse, clicksResponse] = await Promise.all([
      findOwnLinkById(linkId.value),
      getLinkAnalytics(linkId.value, { groupBy: "day" }),
      listRecentClicks(linkId.value, { limit: 20 }),
    ]);

    if (!linkResponse) {
      link.value = null;
      analytics.value = null;
      recentClicks.value = [];
      errorMessage.value = getLinkNotFoundMessage();
      return;
    }

    link.value = linkResponse;
    analytics.value = analyticsResponse;
    recentClicks.value = clicksResponse.items;
  } catch (error: unknown) {
    link.value = null;
    analytics.value = null;
    recentClicks.value = [];
    errorMessage.value = getErrorMessage(error);
  } finally {
    isLoading.value = false;
  }
};

const onLinkUpdated = (updatedLink: Link) => {
  link.value = updatedLink;
};

const onLinkDeleted = async () => {
  await router.push(ROUTES.dashboard);
};

onMounted(() => {
  void loadAnalytics();
});
</script>

<template>
  <section class="link-analytics-page">
    <UiPageHeader
      title="Аналитика ссылки"
      :subtitle="pageSubtitle"
      :back-to="ROUTES.dashboard"
      back-label="Dashboard"
    >
      <template #actions>
        <div v-if="link" class="link-analytics-page__actions">
          <CopyShortUrlButton :short-url="link.shortUrl" variant="secondary" size="md" />
          <UpdateLinkStatusButton :link="link" variant="secondary" size="md" @updated="onLinkUpdated" />
          <DeleteLinkButton
            :link-id="link.id"
            :short-url="link.shortUrl"
            variant="danger"
            size="md"
            @deleted="onLinkDeleted"
          />
        </div>
      </template>
    </UiPageHeader>

    <UiPageState
      v-if="isLoading && !analytics"
      type="loading"
      title="Загружаем аналитику"
      description="Получаем основные метрики и последние переходы."
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

      <section class="link-analytics-page__statistics" aria-labelledby="link-analytics-statistics-title">
        <div class="link-analytics-page__section-header">
          <h2 id="link-analytics-statistics-title" class="link-analytics-page__section-title">Статистика</h2>
        </div>

        <UiPageState
          type="empty"
          title="Графики статистики появятся позже"
          description="Здесь будут графики по переходам, динамике и другим данным ссылки."
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

.link-analytics-page__actions {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  flex-wrap: wrap;
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

.link-analytics-page__statistics,
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

  .link-analytics-page__section-header {
    align-items: stretch;
    flex-direction: column;
  }

  .link-analytics-page__referrer,
  .link-analytics-page__user-agent {
    max-width: 220px;
  }
}
</style>
