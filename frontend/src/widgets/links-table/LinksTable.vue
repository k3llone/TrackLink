<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { Link, LinkStatus, Pagination } from "@/entities/link/link.types";
import { useI18n } from "@/shared/composables/useI18n";
import { getLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiButton, UiPageState, UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";

const router = useRouter();
const { formatDate, formatNumber, t } = useI18n();

const props = withDefaults(
  defineProps<{
    links?: Link[];
    pagination?: Pagination | null;
    loading?: boolean;
    errorMessage?: string;
    hasFilters?: boolean;
  }>(),
  {
    links: () => [],
    pagination: null,
    loading: false,
    errorMessage: "",
    hasFilters: false,
  },
);

const emit = defineEmits<{
  "page-change": [page: number];
  retry: [];
}>();

const columns = computed<UiTableColumn[]>(() => [
  { key: "shortUrl", label: t("common.shortUrl"), width: "22%" },
  { key: "targetUrl", label: t("common.targetUrl"), width: "30%" },
  { key: "createdAt", label: t("common.created"), width: "14%" },
  { key: "status", label: t("common.status"), width: "13%" },
  { key: "totalClicks", label: t("common.clicks"), width: "10%", align: "right" },
]);

const statusLabels: Record<LinkStatus, () => string> = {
  active: () => t("link.status.active"),
  inactive: () => t("link.status.inactive"),
  blocked: () => t("link.status.blocked"),
  deleted: () => t("link.status.deleted"),
};

const currentPage = computed(() => props.pagination?.page ?? 1);
const totalPages = computed(() => props.pagination?.totalPages ?? 0);
const totalItems = computed(() => props.pagination?.totalItems ?? 0);
const hasFilters = computed(() => props.hasFilters);
const showPagination = computed(() => Boolean(props.pagination && totalItems.value > 0));
const canGoPrevious = computed(() => currentPage.value > 1 && !props.loading);
const canGoNext = computed(() => currentPage.value < totalPages.value && !props.loading);
const totalItemsLabel = computed(() => formatNumber(totalItems.value));
const emptyDescription = computed(() =>
  hasFilters.value ? t("links.table.emptyWithFilters") : t("links.table.emptyNoFilters"),
);
const paginationSummary = computed(() =>
  t("links.table.paginationSummary", {
    total: totalItemsLabel.value,
    page: currentPage.value,
    totalPages: totalPages.value,
  }),
);

const getShortUrl = (link: Link) => link.shortUrl || link.code;

const goToPreviousPage = () => {
  if (canGoPrevious.value) {
    emit("page-change", currentPage.value - 1);
  }
};

const goToNextPage = () => {
  if (canGoNext.value) {
    emit("page-change", currentPage.value + 1);
  }
};

const openLinkAnalytics = (row: unknown) => {
  const link = row as Link;
  void router.push(getLinkDetailsPath(link.id));
};

const openLinkAnalyticsById = (linkId: string) => {
  void router.push(getLinkDetailsPath(linkId));
};

const onRetry = () => emit("retry");
</script>

<template>
  <section class="links-table" aria-labelledby="links-table-title">
    <header class="links-table__header">
      <div class="links-table__title-group">
        <h2 id="links-table-title" class="links-table__title">{{ t("links.table.title") }}</h2>
        <p class="links-table__subtitle">{{ t("links.table.subtitle") }}</p>
      </div>
    </header>

    <UiPageState
      v-if="errorMessage"
      type="error"
      :title="t('links.table.errorTitle')"
      :description="errorMessage"
      :action-text="t('common.retry')"
      @action="onRetry"
    />

    <UiTable
      v-else
      :columns="columns"
      :rows="links"
      :loading="loading"
      :empty-text="t('links.table.emptyText')"
      row-clickable
      @row-click="openLinkAnalytics"
    >
      <template #loading>
        <UiPageState
          type="loading"
          :title="t('links.table.loadingTitle')"
          :description="t('links.table.loadingDescription')"
        />
      </template>

      <template #empty>
        <UiPageState
          type="empty"
          :title="t('links.table.emptyTitle')"
          :description="emptyDescription"
        />
      </template>

      <template #cell="{ row, column }">
        <a
          v-if="column.key === 'shortUrl'"
          class="links-table__url links-table__url--short"
          :href="getLinkDetailsPath(row.id)"
          :title="getShortUrl(row)"
          @click.prevent.stop="openLinkAnalyticsById(row.id)"
        >
          {{ getShortUrl(row) }}
        </a>

        <a
          v-else-if="column.key === 'targetUrl'"
          class="links-table__url"
          :href="row.targetUrl"
          target="_blank"
          rel="noreferrer"
          :title="row.targetUrl"
          @click.stop
        >
          {{ row.targetUrl }}
        </a>

        <span v-else-if="column.key === 'createdAt'">{{ formatDate(row.createdAt) }}</span>

        <UiStatusBadge v-else-if="column.key === 'status'" :status="row.status" :label="statusLabels[row.status]()" />

        <span v-else-if="column.key === 'totalClicks'">{{ formatNumber(row.totalClicks) }}</span>

        <span v-else>{{ row[column.key] }}</span>
      </template>
    </UiTable>

    <footer v-if="showPagination" class="links-table__pagination" :aria-label="t('links.table.paginationAria')">
      <span class="links-table__pagination-summary">{{ paginationSummary }}</span>

      <div class="links-table__pagination-actions">
        <UiButton variant="secondary" size="sm" type="button" :disabled="!canGoPrevious" @click="goToPreviousPage">
          {{ t("common.previous") }}
        </UiButton>
        <UiButton variant="secondary" size="sm" type="button" :disabled="!canGoNext" @click="goToNextPage">
          {{ t("common.next") }}
        </UiButton>
      </div>
    </footer>
  </section>
</template>

<style scoped>
.links-table {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.links-table__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.links-table__title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.links-table__title {
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.links-table__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.links-table__url {
  display: inline-block;
  max-width: 260px;
  color: var(--tl-color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: bottom;
  white-space: nowrap;
}

.links-table__url--short {
  color: var(--tl-color-primary);
  font-weight: 700;
}

.links-table__pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.links-table__pagination-actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

@media (max-width: 767px) {
  .links-table__header,
  .links-table__pagination {
    align-items: stretch;
    flex-direction: column;
  }

  .links-table__url {
    max-width: 220px;
  }
}
</style>
