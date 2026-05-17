<script setup lang="ts">
import { computed } from "vue";
import type { AdminLink } from "@/api/admin";
import type { LinkStatus, Pagination } from "@/entities/link/link.types";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton, UiPageState, UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";
import AdminLinkActionsMenu from "./AdminLinkActionsMenu.vue";

const props = withDefaults(
  defineProps<{
    links?: AdminLink[];
    pagination?: Pagination | null;
    loading?: boolean;
    errorMessage?: string;
    hasSearch?: boolean;
  }>(),
  {
    links: () => [],
    pagination: null,
    loading: false,
    errorMessage: "",
    hasSearch: false,
  },
);

const emit = defineEmits<{
  "page-change": [page: number];
  "link-updated": [link: AdminLink];
  retry: [];
}>();

const { formatDate, formatNumber, t } = useI18n();

const columns = computed<UiTableColumn[]>(() => [
  { key: "shortUrl", label: t("common.shortUrl"), width: "18%" },
  { key: "targetUrl", label: t("common.targetUrl"), width: "25%" },
  { key: "owner", label: t("common.owner"), width: "18%" },
  { key: "status", label: t("common.status"), width: "12%" },
  { key: "totalClicks", label: t("common.clicks"), width: "10%", align: "right" },
  { key: "createdAt", label: t("common.created"), width: "12%" },
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
const showPagination = computed(() => Boolean(props.pagination && totalItems.value > 0));
const canGoPrevious = computed(() => currentPage.value > 1 && !props.loading);
const canGoNext = computed(() => currentPage.value < totalPages.value && !props.loading);
const totalItemsLabel = computed(() => formatNumber(totalItems.value));
const emptyStateTitle = computed(() => (props.hasSearch ? t("admin.table.emptySearchTitle") : t("admin.table.emptyDefaultTitle")));
const emptyStateDescription = computed(() =>
  props.hasSearch ? t("admin.table.emptySearchDescription") : t("admin.table.emptyDefaultDescription"),
);
const paginationSummary = computed(() =>
  t("admin.table.paginationSummary", {
    total: totalItemsLabel.value,
    page: currentPage.value,
    totalPages: totalPages.value,
  }),
);

const getShortUrl = (link: AdminLink) => link.shortUrl || link.customAlias || link.code;
const getOwner = (link: AdminLink) => link.ownerEmail?.trim() || link.ownerId;

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

const onRetry = () => emit("retry");
const onLinkUpdated = (link: AdminLink) => emit("link-updated", link);
</script>

<template>
  <section class="admin-links-table" aria-labelledby="admin-links-table-title">
    <header class="admin-links-table__header">
      <div class="admin-links-table__title-group">
        <h2 id="admin-links-table-title" class="admin-links-table__title">{{ t("admin.table.title") }}</h2>
        <p class="admin-links-table__subtitle">{{ t("admin.table.subtitle") }}</p>
      </div>
    </header>

    <UiPageState
      v-if="errorMessage"
      type="error"
      :title="t('admin.table.errorTitle')"
      :description="errorMessage"
      :action-text="t('common.retry')"
      @action="onRetry"
    />

    <UiTable
      v-else
      :columns="columns"
      :rows="links"
      :loading="loading"
      :empty-text="t('admin.table.emptyText')"
    >
      <template #loading>
        <UiPageState
          type="loading"
          :title="t('admin.table.loadingTitle')"
          :description="t('admin.table.loadingDescription')"
        />
      </template>

      <template #empty>
        <UiPageState
          type="empty"
          :title="emptyStateTitle"
          :description="emptyStateDescription"
        />
      </template>

      <template #cell="{ row, column }">
        <span
          v-if="column.key === 'shortUrl'"
          class="admin-links-table__url admin-links-table__url--short"
          :title="getShortUrl(row)"
        >
          {{ getShortUrl(row) }}
        </span>

        <a
          v-else-if="column.key === 'targetUrl'"
          class="admin-links-table__url"
          :href="row.targetUrl"
          target="_blank"
          rel="noreferrer"
          :title="row.targetUrl"
          @click.stop
        >
          {{ row.targetUrl }}
        </a>

        <span v-else-if="column.key === 'owner'" class="admin-links-table__owner" :title="getOwner(row)">
          {{ getOwner(row) }}
        </span>

        <UiStatusBadge v-else-if="column.key === 'status'" :status="row.status" :label="statusLabels[row.status]()" />

        <span v-else-if="column.key === 'totalClicks'">{{ formatNumber(row.totalClicks) }}</span>

        <span v-else-if="column.key === 'createdAt'">{{ formatDate(row.createdAt) }}</span>

        <span v-else>{{ row[column.key] }}</span>
      </template>

      <template #actions="{ row }">
        <AdminLinkActionsMenu :link="row" @blocked="onLinkUpdated" @unblocked="onLinkUpdated" />
      </template>
    </UiTable>

    <footer v-if="showPagination" class="admin-links-table__pagination" :aria-label="t('admin.table.paginationAria')">
      <span class="admin-links-table__pagination-summary">{{ paginationSummary }}</span>

      <div class="admin-links-table__pagination-actions">
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
.admin-links-table {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.admin-links-table__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.admin-links-table__title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.admin-links-table__title {
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.admin-links-table__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.admin-links-table__url,
.admin-links-table__owner {
  display: inline-block;
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: bottom;
  white-space: nowrap;
}

.admin-links-table__url {
  color: var(--tl-color-text);
}

.admin-links-table__url--short {
  color: var(--tl-color-primary);
  font-weight: 700;
}

.admin-links-table__pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.admin-links-table__pagination-actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

@media (max-width: 767px) {
  .admin-links-table__header,
  .admin-links-table__pagination {
    align-items: stretch;
    flex-direction: column;
  }

  .admin-links-table__url,
  .admin-links-table__owner {
    max-width: 220px;
  }
}
</style>
