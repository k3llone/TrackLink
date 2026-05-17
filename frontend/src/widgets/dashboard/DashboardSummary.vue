<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import { UiStatCard } from "@/shared/ui";

const props = withDefaults(
  defineProps<{
    totalLinks?: number;
    activeLinks?: number;
    totalClicks?: number;
    clicksLast24h?: number;
    loading?: boolean;
  }>(),
  {
    totalLinks: 0,
    activeLinks: 0,
    totalClicks: 0,
    clicksLast24h: 0,
    loading: false,
  },
);

const { formatNumber, t } = useI18n();

const cards = computed(() => [
  {
    key: "totalLinks",
    title: t("dashboard.summary.totalLinks.title"),
    value: formatNumber(props.totalLinks),
    hint: t("dashboard.summary.totalLinks.hint"),
  },
  {
    key: "activeLinks",
    title: t("dashboard.summary.activeLinks.title"),
    value: formatNumber(props.activeLinks),
    hint: t("dashboard.summary.activeLinks.hint"),
  },
  {
    key: "totalClicks",
    title: t("dashboard.summary.totalClicks.title"),
    value: formatNumber(props.totalClicks),
    hint: t("dashboard.summary.totalClicks.hint"),
  },
  {
    key: "clicksLast24h",
    title: t("dashboard.summary.clicksLast24h.title"),
    value: formatNumber(props.clicksLast24h),
    hint: t("dashboard.summary.clicksLast24h.hint"),
  },
]);
</script>

<template>
  <section class="dashboard-summary" :aria-label="t('dashboard.summary.aria')">
    <UiStatCard
      v-for="card in cards"
      :key="card.key"
      :title="card.title"
      :value="card.value"
      :hint="card.hint"
      :loading="loading"
    />
  </section>
</template>

<style scoped>
.dashboard-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 28px;
}

@media (max-width: 1023px) {
  .dashboard-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .dashboard-summary {
    grid-template-columns: 1fr;
    margin-bottom: 24px;
  }
}
</style>
