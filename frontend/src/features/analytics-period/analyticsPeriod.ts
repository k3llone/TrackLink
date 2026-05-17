import type { AnalyticsGroupBy } from "@/entities/analytics/analytics.types";

const HOUR_MS = 60 * 60 * 1000;
const DAY_MS = 24 * HOUR_MS;

export type AnalyticsPeriodValue = "24h" | "7d" | "30d";

export interface AnalyticsPeriodOption {
  value: AnalyticsPeriodValue;
  label: string;
  groupBy: AnalyticsGroupBy;
  durationMs: number;
}

export interface AnalyticsPeriodParams {
  from: string;
  to: string;
  groupBy: AnalyticsGroupBy;
}

export const DEFAULT_ANALYTICS_PERIOD: AnalyticsPeriodValue = "7d";

export const ANALYTICS_PERIOD_OPTIONS: AnalyticsPeriodOption[] = [
  {
    value: "24h",
    label: "24h",
    groupBy: "hour",
    durationMs: 24 * HOUR_MS,
  },
  {
    value: "7d",
    label: "7d",
    groupBy: "day",
    durationMs: 7 * DAY_MS,
  },
  {
    value: "30d",
    label: "30d",
    groupBy: "day",
    durationMs: 30 * DAY_MS,
  },
];

export const getAnalyticsPeriodOption = (value: AnalyticsPeriodValue) =>
  ANALYTICS_PERIOD_OPTIONS.find((option) => option.value === value) ?? ANALYTICS_PERIOD_OPTIONS[1];

export const getAnalyticsPeriodParams = (
  value: AnalyticsPeriodValue,
  now = new Date(),
): AnalyticsPeriodParams => {
  const option = getAnalyticsPeriodOption(value);
  const to = new Date(now);
  const from = new Date(to.getTime() - option.durationMs);

  return {
    from: from.toISOString(),
    to: to.toISOString(),
    groupBy: option.groupBy,
  };
};
