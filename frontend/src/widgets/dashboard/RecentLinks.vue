<script setup lang="ts">
import { useRouter } from "vue-router";
import type { Link } from "@/entities/link/link.types";
import { getLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";

defineProps<{
  links: Link[];
}>();

const router = useRouter();

const columns: UiTableColumn[] = [
  { key: "shortUrl", label: "Short URL", width: "24%" },
  { key: "targetUrl", label: "Target URL", width: "30%" },
  { key: "createdAt", label: "Создана", width: "16%" },
  { key: "status", label: "Статус", width: "14%" },
  { key: "totalClicks", label: "Переходы", width: "12%", align: "right" },
];

const dateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
  year: "numeric",
});

const numberFormatter = new Intl.NumberFormat("ru-RU");

const statusLabels: Record<Link["status"], string> = {
  active: "Активна",
  inactive: "Неактивна",
  blocked: "Заблокирована",
  deleted: "Удалена",
};

const formatDate = (value: string) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "—";
  }

  return dateFormatter.format(date);
};

const formatNumber = (value: number) => numberFormatter.format(value);

const openLinkAnalytics = (row: unknown) => {
  const link = row as Link;
  void router.push(getLinkDetailsPath(link.id));
};

const openLinkAnalyticsById = (linkId: string) => {
  void router.push(getLinkDetailsPath(linkId));
};
</script>

<template>
  <section class="recent-links" aria-labelledby="recent-links-title">
    <header class="recent-links__header">
      <div class="recent-links__title-group">
        <h2 id="recent-links-title" class="recent-links__title">Последние ссылки</h2>
        <p class="recent-links__subtitle">Компактный список недавно созданных ссылок аккаунта.</p>
      </div>
    </header>

    <UiTable :columns="columns" :rows="links" empty-text="Последних ссылок пока нет." row-clickable @row-click="openLinkAnalytics">
      <template #cell="{ row, column }">
        <a
          v-if="column.key === 'shortUrl'"
          class="recent-links__url recent-links__url--short"
          :href="getLinkDetailsPath(row.id)"
          :title="row.shortUrl"
          @click.prevent.stop="openLinkAnalyticsById(row.id)"
        >
          {{ row.shortUrl }}
        </a>

        <a
          v-else-if="column.key === 'targetUrl'"
          class="recent-links__url"
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
    </UiTable>
  </section>
</template>

<style scoped>
.recent-links {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.recent-links__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.recent-links__title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.recent-links__title {
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.recent-links__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.recent-links__url {
  display: inline-block;
  max-width: 280px;
  color: var(--tl-color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: bottom;
  white-space: nowrap;
}

.recent-links__url--short {
  color: var(--tl-color-primary);
  font-weight: 700;
}

@media (max-width: 767px) {
  .recent-links__url {
    max-width: 220px;
  }
}
</style>
