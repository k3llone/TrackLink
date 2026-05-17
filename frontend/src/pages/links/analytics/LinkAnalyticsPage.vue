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
import {
  AnalyticsPeriodPicker,
  DEFAULT_ANALYTICS_PERIOD,
  getAnalyticsPeriodOption,
  getAnalyticsPeriodParams,
  type AnalyticsPeriodValue,
} from "@/features/analytics-period";
import CopyShortUrlButton from "@/features/link-actions/CopyShortUrlButton.vue";
import DeleteLinkButton from "@/features/link-actions/DeleteLinkButton.vue";
import UpdateLinkStatusButton from "@/features/link-actions/UpdateLinkStatusButton.vue";
import { useI18n } from "@/shared/composables/useI18n";
import { ROUTES } from "@/shared/lib/routes/paths";
import {
  UiPageHeader,
  UiPageState,
  UiStatCard,
  UiTable,
  type UiTableColumn,
} from "@/shared/ui";
import { ClicksTimeChart } from "@/widgets/analytics-chart";

const route = useRoute();
const router = useRouter();
const { formatDateTime, formatNumber, t } = useI18n();

const linkId = computed(() => {
  const rawId = route.params.id;
  return Array.isArray(rawId) ? rawId[0] : rawId;
});

const analytics = ref<LinkAnalyticsResponse | null>(null);
const link = ref<Link | null>(null);
const recentClicks = ref<ClickEvent[]>([]);
const isLoading = ref(false);
const isAnalyticsLoading = ref(false);
const errorMessage = ref("");
const analyticsErrorMessage = ref("");
const selectedPeriod = ref<AnalyticsPeriodValue>(DEFAULT_ANALYTICS_PERIOD);
let analyticsRequestId = 0;

const clicksColumns = computed<UiTableColumn[]>(() => [
  { key: "clickedAt", label: t("analytics.table.clickedAt"), width: "24%" },
  { key: "referrer", label: t("analytics.table.referrer"), width: "30%" },
  { key: "userAgent", label: t("common.userAgent"), width: "46%" },
]);

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const currentGroupBy = computed(() => getAnalyticsPeriodOption(selectedPeriod.value).groupBy);
const totalClicks = computed(() => formatNumber(analytics.value?.totalClicks ?? link.value?.totalClicks ?? 0));
const clicksLast24h = computed(() => formatNumber(analytics.value?.clicksLast24h ?? 0));
const pageSubtitle = computed(() => {
  if (link.value) {
    return t("analytics.page.subtitleWithUrl", { shortUrl: link.value.shortUrl });
  }

  return linkId.value
    ? t("analytics.page.subtitleWithId", { linkId: linkId.value })
    : t("analytics.page.subtitleNoLink");
});

const lastClickedAt = computed(() => {
  const value = analytics.value?.lastClickedAt ?? link.value?.lastClickedAt;

  if (!value) {
    return t("analytics.metrics.noClicks");
  }

  return formatDateTime(value);
});

const formatSource = (value?: string | null) => {
  const source = value?.trim();
  return source || t("common.directVisit");
};

const formatUserAgent = (value?: string | null) => {
  const userAgent = value?.trim();
  return userAgent || t("common.unknown");
};

const getErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401) {
      return t("analytics.error.unauthorized");
    }

    if (error.status === 404) {
      return t("analytics.error.notFound");
    }
  }

  return t("analytics.error.failed");
};

const getLinkNotFoundMessage = () => t("analytics.notFound");

const loadAnalyticsSeries = async () => {
  if (!linkId.value || isAnalyticsLoading.value) {
    return;
  }

  const requestId = ++analyticsRequestId;
  isAnalyticsLoading.value = true;
  analyticsErrorMessage.value = "";

  try {
    const response = await getLinkAnalytics(linkId.value, getAnalyticsPeriodParams(selectedPeriod.value));

    if (requestId !== analyticsRequestId) {
      return;
    }

    analytics.value = response;
  } catch (error: unknown) {
    if (requestId !== analyticsRequestId) {
      return;
    }

    analytics.value = null;
    analyticsErrorMessage.value = getErrorMessage(error);
  } finally {
    if (requestId === analyticsRequestId) {
      isAnalyticsLoading.value = false;
    }
  }
};

const loadAnalytics = async () => {
  if (!linkId.value || isLoading.value) {
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";
  analyticsErrorMessage.value = "";

  try {
    const [linkResponse, clicksResponse] = await Promise.all([
      findOwnLinkById(linkId.value),
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
    recentClicks.value = clicksResponse.items;
    await loadAnalyticsSeries();
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

const onAnalyticsPeriodChange = (period: AnalyticsPeriodValue) => {
  selectedPeriod.value = period;
  void loadAnalyticsSeries();
};

onMounted(() => {
  void loadAnalytics();
});
</script>

<template>
  <section class="link-analytics-page">
    <UiPageHeader
      :title="t('analytics.page.title')"
      :subtitle="pageSubtitle"
      :back-to="ROUTES.dashboard"
      :back-label="t('common.dashboard')"
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
      v-if="isLoading && !link"
      type="loading"
      :title="t('analytics.loading.title')"
      :description="t('analytics.loading.description')"
    />

    <UiPageState
      v-else-if="errorMessage"
      type="error"
      :title="t('analytics.unavailable.title')"
      :description="errorMessage"
      :action-text="t('common.retry')"
      @action="loadAnalytics"
    />

    <UiPageState
      v-else-if="!linkId"
      type="not-found"
      :title="t('analytics.noLink.title')"
      :description="t('analytics.noLink.description')"
      :action-to="ROUTES.dashboard"
      :action-text="t('analytics.noLink.action')"
    />

    <div v-else-if="link" class="link-analytics-page__content">
      <section class="link-analytics-page__summary" :aria-label="t('analytics.metrics.aria')">
        <UiStatCard :title="t('analytics.metrics.totalClicks')" :value="totalClicks" />
        <UiStatCard :title="t('analytics.metrics.clicksLast24h')" :value="clicksLast24h" />
        <UiStatCard :title="t('analytics.metrics.lastClick')" :value="lastClickedAt" />
      </section>

      <section class="link-analytics-page__statistics" aria-labelledby="link-analytics-statistics-title">
        <div class="link-analytics-page__section-header">
          <h2 id="link-analytics-statistics-title" class="link-analytics-page__section-title">{{ t("analytics.statistics.title") }}</h2>
          <AnalyticsPeriodPicker
            v-model="selectedPeriod"
            :loading="isAnalyticsLoading"
            @change="onAnalyticsPeriodChange"
          />
        </div>

        <ClicksTimeChart
          :series="analytics?.series ?? []"
          :group-by="currentGroupBy"
          :loading="isAnalyticsLoading"
          :error="analyticsErrorMessage"
          @retry="loadAnalyticsSeries"
        />
      </section>

      <section class="link-analytics-page__clicks" aria-labelledby="link-analytics-clicks-title">
        <div class="link-analytics-page__section-header">
          <h2 id="link-analytics-clicks-title" class="link-analytics-page__section-title">{{ t("analytics.recentClicks.title") }}</h2>
        </div>

        <UiTable :columns="clicksColumns" :rows="recentClicks" :empty-text="t('analytics.recentClicks.emptyText')">
          <template #empty>
            <UiPageState
              type="empty"
              :title="t('analytics.recentClicks.emptyTitle')"
              :description="t('analytics.recentClicks.emptyDescription')"
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
