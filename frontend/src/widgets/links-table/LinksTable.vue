<script setup lang="ts">
import { computed } from "vue";
import type { Link, LinkStatus, Pagination } from "@/entities/link/link.types";
import { UiButton, UiInput, UiPageState, UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";
import LinkRowActions from "./LinkRowActions.vue";

const props = withDefaults(
  defineProps<{
    links?: Link[];
    pagination?: Pagination | null;
    loading?: boolean;
    errorMessage?: string;
    q?: string;
    status?: LinkStatus | "";
  }>(),
  {
    links: () => [],
    pagination: null,
    loading: false,
    errorMessage: "",
    q: "",
    status: "",
  },
);

const emit = defineEmits<{
  "filters-change": [filters: { q: string; status: LinkStatus | "" }];
  "page-change": [page: number];
  retry: [];
}>();

const columns: UiTableColumn[] = [
  { key: "shortUrl", label: "Short URL", width: "22%" },
  { key: "targetUrl", label: "Target URL", width: "30%" },
  { key: "createdAt", label: "Создана", width: "14%" },
  { key: "status", label: "Статус", width: "13%" },
  { key: "totalClicks", label: "Переходы", width: "10%", align: "right" },
];

const statusOptions: Array<{ value: LinkStatus | ""; label: string }> = [
  { value: "", label: "Все статусы" },
  { value: "active", label: "Активные" },
  { value: "inactive", label: "Неактивные" },
  { value: "blocked", label: "Заблокированные" },
  { value: "deleted", label: "Удаленные" },
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
const hasFilters = computed(() => Boolean(props.q.trim() || props.status));
const showPagination = computed(() => Boolean(props.pagination && totalItems.value > 0));
const canGoPrevious = computed(() => currentPage.value > 1 && !props.loading);
const canGoNext = computed(() => currentPage.value < totalPages.value && !props.loading);
const totalItemsLabel = computed(() => numberFormatter.format(totalItems.value));
const emptyDescription = computed(() =>
  hasFilters.value
    ? "По текущим фильтрам ссылок не найдено."
    : "Создайте первую короткую ссылку, чтобы она появилась в списке.",
);

const formatDate = (value: string) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "—";
  }

  return dateFormatter.format(date);
};

const formatNumber = (value: number) => numberFormatter.format(value);

const getShortUrl = (link: Link) => link.shortUrl || link.code;

const onSearchChange = (q: string) => {
  emit("filters-change", { q, status: props.status });
};

const onStatusChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  emit("filters-change", { q: props.q, status: target.value as LinkStatus | "" });
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
</script>

<template>
  <section class="links-table" aria-labelledby="links-table-title">
    <header class="links-table__header">
      <div class="links-table__title-group">
        <h2 id="links-table-title" class="links-table__title">Ссылки</h2>
        <p class="links-table__subtitle">Список коротких ссылок аккаунта.</p>
      </div>

      <div class="links-table__filters" aria-label="Фильтры списка ссылок">
        <UiInput
          :model-value="q"
          class="links-table__search"
          type="search"
          placeholder="Поиск по ссылке или URL"
          autocomplete="off"
          @update:model-value="onSearchChange"
        />

        <label class="links-table__status-filter">
          <span class="links-table__status-label">Статус</span>
          <select class="links-table__status-select" :value="status" @change="onStatusChange">
            <option v-for="option in statusOptions" :key="option.value || 'all'" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </label>
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

    <UiTable v-else :columns="columns" :rows="links" :loading="loading" empty-text="Ссылок пока нет.">
      <template #loading>
        <UiPageState
          type="loading"
          title="Загружаем ссылки"
          description="Получаем список коротких ссылок аккаунта."
        />
      </template>

      <template #empty>
        <UiPageState
          type="empty"
          title="Ссылок пока нет"
          :description="emptyDescription"
        />
      </template>

      <template #cell="{ row, column }">
        <a
          v-if="column.key === 'shortUrl'"
          class="links-table__url links-table__url--short"
          :href="getShortUrl(row)"
          target="_blank"
          rel="noreferrer"
          :title="getShortUrl(row)"
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
        >
          {{ row.targetUrl }}
        </a>

        <span v-else-if="column.key === 'createdAt'">{{ formatDate(row.createdAt) }}</span>

        <UiStatusBadge v-else-if="column.key === 'status'" :status="row.status" :label="statusLabels[row.status]" />

        <span v-else-if="column.key === 'totalClicks'">{{ formatNumber(row.totalClicks) }}</span>

        <span v-else>{{ row[column.key] }}</span>
      </template>

      <template #actions="{ row }">
        <LinkRowActions :link="row" />
      </template>
    </UiTable>

    <footer v-if="showPagination" class="links-table__pagination" aria-label="Пагинация списка ссылок">
      <span class="links-table__pagination-summary">
        {{ totalItemsLabel }} ссылок · страница {{ currentPage }} из {{ totalPages }}
      </span>

      <div class="links-table__pagination-actions">
        <UiButton variant="secondary" size="sm" type="button" :disabled="!canGoPrevious" @click="goToPreviousPage">
          Назад
        </UiButton>
        <UiButton variant="secondary" size="sm" type="button" :disabled="!canGoNext" @click="goToNextPage">
          Вперед
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

.links-table__filters {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  width: min(100%, 540px);
}

.links-table__search {
  min-width: 260px;
}

.links-table__status-filter {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 180px;
  color: var(--tl-color-text);
}

.links-table__status-label {
  font-size: 13px;
  font-weight: 700;
}

.links-table__status-select {
  min-height: 44px;
  border: 1px solid #ddd7e8;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
  color: var(--tl-color-text);
  font-family: var(--tl-font-family);
  font-size: 14px;
  padding: 0 12px;
}

.links-table__status-select:focus {
  border-color: var(--tl-color-primary);
  box-shadow: 0 0 0 2px rgb(109 74 255 / 18%);
  outline: 0;
}

.links-table__status-select:disabled {
  opacity: 0.65;
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
  .links-table__filters,
  .links-table__pagination {
    align-items: stretch;
    flex-direction: column;
  }

  .links-table__filters {
    width: 100%;
  }

  .links-table__search {
    min-width: 0;
  }

  .links-table__url {
    max-width: 220px;
  }
}
</style>
