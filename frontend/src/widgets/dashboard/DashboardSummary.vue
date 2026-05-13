<script setup lang="ts">
import { computed } from "vue";
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

const numberFormatter = new Intl.NumberFormat("ru-RU");

const formatNumber = (value: number) => numberFormatter.format(value);

const cards = computed(() => [
  {
    key: "totalLinks",
    title: "Всего ссылок",
    value: formatNumber(props.totalLinks),
    hint: "Созданные короткие ссылки",
  },
  {
    key: "activeLinks",
    title: "Активные ссылки",
    value: formatNumber(props.activeLinks),
    hint: "Сейчас доступны для переходов",
  },
  {
    key: "totalClicks",
    title: "Всего переходов",
    value: formatNumber(props.totalClicks),
    hint: "Суммарно по всем ссылкам",
  },
  {
    key: "clicksLast24h",
    title: "За 24 часа",
    value: formatNumber(props.clicksLast24h),
    hint: "Переходы за последние сутки",
  },
]);
</script>

<template>
  <section class="dashboard-summary" aria-label="Ключевые показатели dashboard">
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
