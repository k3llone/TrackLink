<script setup lang="ts">
import type { Link, LinkStatus } from "@/entities/link/link.types";
import { UiInput, UiPageState, UiTable, type UiTableColumn } from "@/shared/ui";

const props = withDefaults(
  defineProps<{
    links?: Link[];
    loading?: boolean;
    errorMessage?: string;
    q?: string;
    status?: LinkStatus | "";
  }>(),
  {
    links: () => [],
    loading: false,
    errorMessage: "",
    q: "",
    status: "",
  },
);

const emit = defineEmits<{
  "filters-change": [filters: { q: string; status: LinkStatus | "" }];
  retry: [];
}>();

const columns: UiTableColumn[] = [
  { key: "shortUrl", label: "Short URL", width: "24%" },
  { key: "targetUrl", label: "Target URL", width: "34%" },
  { key: "createdAt", label: "Создана", width: "16%" },
  { key: "status", label: "Статус", width: "14%" },
  { key: "totalClicks", label: "Переходы", width: "12%", align: "right" },
];

const statusOptions: Array<{ value: LinkStatus | ""; label: string }> = [
  { value: "", label: "Все статусы" },
  { value: "active", label: "Активные" },
  { value: "inactive", label: "Неактивные" },
  { value: "blocked", label: "Заблокированные" },
  { value: "deleted", label: "Удаленные" },
];

const onSearchChange = (q: string) => {
  emit("filters-change", { q, status: props.status });
};

const onStatusChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  emit("filters-change", { q: props.q, status: target.value as LinkStatus | "" });
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
          :disabled="loading"
          @update:model-value="onSearchChange"
        />

        <label class="links-table__status-filter">
          <span class="links-table__status-label">Статус</span>
          <select class="links-table__status-select" :value="status" :disabled="loading" @change="onStatusChange">
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
          description="Создайте первую короткую ссылку, чтобы она появилась в списке."
        />
      </template>
    </UiTable>
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

@media (max-width: 767px) {
  .links-table__header,
  .links-table__filters {
    align-items: stretch;
    flex-direction: column;
  }

  .links-table__filters {
    width: 100%;
  }

  .links-table__search {
    min-width: 0;
  }
}
</style>
