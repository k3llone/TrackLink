<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { Link } from "@/entities/link/link.types";
import { useI18n } from "@/shared/composables/useI18n";
import { getLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";

defineProps<{
  links: Link[];
}>();

const router = useRouter();
const { formatDate, formatNumber, t } = useI18n();

const columns = computed<UiTableColumn[]>(() => [
  { key: "shortUrl", label: t("common.shortUrl"), width: "24%" },
  { key: "targetUrl", label: t("common.targetUrl"), width: "30%" },
  { key: "createdAt", label: t("common.created"), width: "16%" },
  { key: "status", label: t("common.status"), width: "14%" },
  { key: "totalClicks", label: t("common.clicks"), width: "12%", align: "right" },
]);

const statusLabels: Record<Link["status"], () => string> = {
  active: () => t("link.status.active"),
  inactive: () => t("link.status.inactive"),
  blocked: () => t("link.status.blocked"),
  deleted: () => t("link.status.deleted"),
};

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
        <h2 id="recent-links-title" class="recent-links__title">{{ t("links.recent.title") }}</h2>
        <p class="recent-links__subtitle">{{ t("links.recent.subtitle") }}</p>
      </div>
    </header>

    <UiTable :columns="columns" :rows="links" :empty-text="t('links.recent.emptyText')" row-clickable @row-click="openLinkAnalytics">
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

        <UiStatusBadge v-else-if="column.key === 'status'" :status="row.status" :label="statusLabels[row.status]()" />

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
