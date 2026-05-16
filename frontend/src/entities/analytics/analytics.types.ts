export type AnalyticsGroupBy = "hour" | "day";

export interface TimeSeriesPoint {
  periodStart: string;
  clicks: number;
}
