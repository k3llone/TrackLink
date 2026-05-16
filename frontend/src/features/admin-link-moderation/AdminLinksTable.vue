<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { AdminLink } from "@/api/admin";
import type { LinkStatus, Pagination } from "@/entities/link/link.types";
import { getAdminLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiButton, UiPageState, UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";
import AdminLinkActionsMenu from "./AdminLinkActionsMenu.vue";

const router = useRouter();

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

const columns: UiTableColumn[] = [
  { key: "shortUrl", label: "Short URL", width: "18%" },
  { key: "targetUrl", label: "Target URL", width: "25%" },
  { key: "owner", label: "Owner", width: "18%" },
  { key: "status", label: "Статус", width: "12%" },
  { key: "totalClicks", label: "Переходы", width: "10%", align: "right" },
  { key: "createdAt", label: "Создана", width: "12%" },
];

const dateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
  year: "numeric",
});

const numberFormatter = new Intl.NumberFormat("ru-RU");

const statusLabels: Record<LinkStatus, string> = {
  active: "Активна",
  inactive: "Неактивна",
  blocked: "Заблокирована",
  deleted: "Удалена",
};

const currentPage = computed(() => props.pagination?.page ?? 1);
const totalPages = computed(() => props.pagination?.totalPages ?? 0);
const totalItems = computed(() => props.pagination?.totalItems ?? 0);
const showPagination = computed(() => Boolean(props.pagination && totalItems.value > 0));
const canGoPrevious = computed(() => currentPage.value > 1 && !props.loading);
const canGoNext = computed(() => currentPage.value < totalPages.value && !props.loading);
const totalItemsLabel = computed(() => numberFormatter.format(totalItems.value));
const emptyStateTitle = computed(() => (props.hasSearch ? "Ссылки не найдены" : "Ссылок пока нет"));
const emptyStateDescription = computed(() =>
  props.hasSearch ? "По запросу ничего не найдено." : "Административный список ссылок пуст.",
);

const formatDate = (value: string) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "—";
  }

  return dateFormatter.format(date);
};

const formatNumber = (value: number) => numberFormatter.format(value);

const getShortUrl = (link: AdminLink) => link.shortUrl || link.customAlias || link.code;
const getOwner = (link: AdminLink) => link.ownerEmail?.trim() || link.ownerId;

const openLinkAnalyticsById = (linkId: string) => {
  void router.push(getAdminLinkDetailsPath(linkId));
};

const openLinkAnalytics = (row: unknown) => {
  const link = row as AdminLink;
  openLinkAnalyticsById(link.id);
};

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
        <h2 id="admin-links-table-title" class="admin-links-table__title">Ссылки</h2>
        <p class="admin-links-table__subtitle">
          Административный список коротких ссылок для просмотра и блокировки.
        </p>
      </div>
    </header>

    <UiPageState
      v-if="errorMessage"
      type="error"
      title="Список ссылок недоступен"
      :description="errorMessage"
      action-text="Повторить"
      @action="onRetry"
    />

    <UiTable
      v-else
      :columns="columns"
      :rows="links"
      :loading="loading"
      empty-text="Ссылок пока нет."
      row-clickable
      @row-click="openLinkAnalytics"
    >
      <template #loading>
        <UiPageState
          type="loading"
          title="Загружаем ссылки"
          description="Получаем административный список коротких ссылок."
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
        <a
          v-if="column.key === 'shortUrl'"
          class="admin-links-table__url admin-links-table__url--short"
          :href="getAdminLinkDetailsPath(row.id)"
          :title="getShortUrl(row)"
          @click.prevent.stop="openLinkAnalyticsById(row.id)"
        >
          {{ getShortUrl(row) }}
        </a>

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

        <UiStatusBadge v-else-if="column.key === 'status'" :status="row.status" :label="statusLabels[row.status]" />

        <span v-else-if="column.key === 'totalClicks'">{{ formatNumber(row.totalClicks) }}</span>

        <span v-else-if="column.key === 'createdAt'">{{ formatDate(row.createdAt) }}</span>

        <span v-else>{{ row[column.key] }}</span>
      </template>

      <template #actions="{ row }">
        <AdminLinkActionsMenu :link="row" @blocked="onLinkUpdated" />
      </template>
    </UiTable>

    <footer v-if="showPagination" class="admin-links-table__pagination" aria-label="Пагинация списка ссылок">
      <span class="admin-links-table__pagination-summary">
        {{ totalItemsLabel }} ссылок · страница {{ currentPage }} из {{ totalPages }}
      </span>

      <div class="admin-links-table__pagination-actions">
        <UiButton variant="secondary" size="sm" type="button" :disabled="!canGoPrevious" @click="goToPreviousPage">
          Назад
        </UiButton>
        <UiButton variant="secondary" size="sm" type="button" :disabled="!canGoNext" @click="goToNextPage">
          Вперёд
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
