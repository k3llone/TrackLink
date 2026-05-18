<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { Link } from "@/entities/link/link.types";
import { useCopyShortUrl } from "@/features/link-actions/useCopyShortUrl";
import { useI18n } from "@/shared/composables/useI18n";
import { getLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiStatusBadge, UiTable, type UiTableColumn } from "@/shared/ui";

defineProps<{
  links: Link[];
}>();

const router = useRouter();
const { formatDate, formatNumber, t } = useI18n();
const { copyShortUrl } = useCopyShortUrl();

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

const getShortUrl = (link: Link) => link.shortUrl || link.code;

const openLinkAnalytics = (row: unknown) => {
  const link = row as Link;
  void router.push(getLinkDetailsPath(link.id));
};

const getCopyShortUrlLabel = (link: Link) =>
  t("linkActions.copy.shortUrlTitle", { shortUrl: getShortUrl(link) });

const copyLinkShortUrl = (link: Link) => {
  void copyShortUrl(getShortUrl(link));
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
        <button
          v-if="column.key === 'shortUrl'"
          class="recent-links__url recent-links__url--short recent-links__url--copy"
          type="button"
          :aria-label="getCopyShortUrlLabel(row)"
          :title="getCopyShortUrlLabel(row)"
          @click.stop="copyLinkShortUrl(row)"
          @keydown.enter.stop
          @keydown.space.stop
        >
          {{ getShortUrl(row) }}
        </button>

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

.recent-links__url--copy {
  border: 0;
  padding: 0;
  background: transparent;
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.recent-links__url--copy:hover,
.recent-links__url--copy:focus-visible {
  text-decoration: underline;
}

.recent-links__url--copy:focus-visible {
  outline: 2px solid rgb(109 74 255 / 35%);
  outline-offset: 2px;
}

@media (max-width: 767px) {
  .recent-links__url {
    max-width: 220px;
  }
}
</style>
