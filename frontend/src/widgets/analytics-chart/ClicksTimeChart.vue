<script setup lang="ts">
import { computed } from "vue";
import type { AnalyticsGroupBy, TimeSeriesPoint } from "@/entities/analytics/analytics.types";
import { useI18n } from "@/shared/composables/useI18n";
import { UiPageState } from "@/shared/ui";
import { buildClicksChartData } from "./chartData";

const props = withDefaults(
  defineProps<{
    series: TimeSeriesPoint[];
    groupBy: AnalyticsGroupBy;
    loading?: boolean;
    error?: string | null;
  }>(),
  {
    loading: false,
    error: null,
  },
);

const emit = defineEmits<{
  retry: [];
}>();

const { t } = useI18n();
const safeSeries = computed(() => props.series ?? []);
const chartData = computed(() => buildClicksChartData(safeSeries.value, props.groupBy));
const hasError = computed(() => Boolean(props.error?.trim()));
const hasClicks = computed(() => safeSeries.value.some((point) => Number.isFinite(point.clicks) && point.clicks > 0));
const isEmpty = computed(() => !safeSeries.value.length || !hasClicks.value);
const periodLabel = computed(() => (props.groupBy === "hour" ? t("analytics.chart.periodHours") : t("analytics.chart.periodDays")));
const ariaLabel = computed(() => t("analytics.chart.aria", { period: periodLabel.value }));

const onRetry = () => emit("retry");
</script>

<template>
  <section class="clicks-time-chart" aria-labelledby="clicks-time-chart-title">
    <header class="clicks-time-chart__header">
      <div class="clicks-time-chart__title-group">
        <h3 id="clicks-time-chart-title" class="clicks-time-chart__title">{{ t("analytics.chart.title") }}</h3>
        <p class="clicks-time-chart__subtitle">{{ t("analytics.chart.subtitle", { period: periodLabel }) }}</p>
      </div>
    </header>

    <UiPageState
      v-if="loading"
      type="loading"
      :title="t('analytics.chart.loadingTitle')"
      :description="t('analytics.chart.loadingDescription')"
    />

    <UiPageState
      v-else-if="hasError"
      type="error"
      :title="t('analytics.chart.errorTitle')"
      :description="error || t('analytics.chart.errorFallback')"
      :action-text="t('common.retry')"
      @action="onRetry"
    />

    <UiPageState
      v-else-if="isEmpty"
      type="empty"
      :title="t('analytics.chart.emptyTitle')"
      :description="t('analytics.chart.emptyDescription')"
    />

    <div v-else class="clicks-time-chart__scroll">
      <div class="clicks-time-chart__canvas" role="img" :aria-label="ariaLabel">
        <svg
          class="clicks-time-chart__svg"
          :viewBox="`0 0 ${chartData.width} ${chartData.height}`"
          :width="chartData.width"
          :height="chartData.height"
          focusable="false"
          aria-hidden="true"
        >
          <g class="clicks-time-chart__grid">
            <g v-for="label in chartData.yAxisLabels" :key="label.key">
              <line x1="56" x2="734" :y1="label.y" :y2="label.y" />
              <text x="44" :y="label.y" dy="4" text-anchor="end">{{ label.label }}</text>
            </g>
          </g>

          <path v-if="chartData.areaPath" class="clicks-time-chart__area" :d="chartData.areaPath" />
          <path v-if="chartData.points.length > 1" class="clicks-time-chart__line" :d="chartData.linePath" />

          <g class="clicks-time-chart__points">
            <circle
              v-for="point in chartData.points"
              :key="point.key"
              class="clicks-time-chart__point"
              :cx="point.x"
              :cy="point.y"
              r="4"
            >
              <title>{{ point.tooltip }}</title>
            </circle>
          </g>

          <g class="clicks-time-chart__x-axis">
            <text
              v-for="label in chartData.xAxisLabels"
              :key="label.key"
              :x="label.x"
              y="276"
              text-anchor="middle"
            >
              {{ label.label }}
            </text>
          </g>
        </svg>
      </div>
    </div>
  </section>
</template>

<style scoped>
.clicks-time-chart {
  display: flex;
  flex-direction: column;
  gap: 14px;
  width: 100%;
  padding: 18px;
  border: 1px solid #ece7f4;
  border-radius: var(--tl-radius-lg);
  background: var(--tl-color-white);
}

.clicks-time-chart__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.clicks-time-chart__title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.clicks-time-chart__title {
  margin: 0;
  color: var(--tl-color-text);
  font-size: 18px;
  line-height: 1.25;
}

.clicks-time-chart__subtitle {
  margin: 0;
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.clicks-time-chart__scroll {
  width: 100%;
  overflow-x: auto;
}

.clicks-time-chart__canvas {
  min-width: 680px;
}

.clicks-time-chart__svg {
  display: block;
  width: 100%;
  height: auto;
}

.clicks-time-chart__grid line {
  stroke: rgb(37 31 63 / 9%);
  stroke-width: 1;
}

.clicks-time-chart__grid text,
.clicks-time-chart__x-axis text {
  fill: var(--tl-color-text-muted);
  font-size: 12px;
}

.clicks-time-chart__area {
  fill: rgb(109 74 255 / 10%);
}

.clicks-time-chart__line {
  fill: none;
  stroke: var(--tl-color-primary);
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 3;
}

.clicks-time-chart__point {
  fill: var(--tl-color-white);
  stroke: var(--tl-color-primary);
  stroke-width: 3;
}

@media (max-width: 767px) {
  .clicks-time-chart {
    padding: 16px;
  }

  .clicks-time-chart__canvas {
    min-width: 640px;
  }
}
</style>
