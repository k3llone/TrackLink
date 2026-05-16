import type { AnalyticsGroupBy, TimeSeriesPoint } from "@/entities/analytics/analytics.types";

const CHART_WIDTH = 760;
const CHART_HEIGHT = 300;
const PADDING = {
  top: 24,
  right: 26,
  bottom: 48,
  left: 56,
};
const MAX_X_LABELS = 6;

const numberFormatter = new Intl.NumberFormat("ru-RU");
const hourFormatter = new Intl.DateTimeFormat("ru-RU", {
  hour: "2-digit",
  minute: "2-digit",
});
const dayFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
});
const fullDateFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "short",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});

export interface ClicksChartPoint {
  key: string;
  x: number;
  y: number;
  clicks: number;
  xLabel: string;
  tooltip: string;
}

export interface ChartAxisLabel {
  key: string;
  label: string;
  x?: number;
  y?: number;
}

export interface ClicksChartData {
  width: number;
  height: number;
  linePath: string;
  areaPath: string;
  points: ClicksChartPoint[];
  xAxisLabels: ChartAxisLabel[];
  yAxisLabels: ChartAxisLabel[];
}

const parseDate = (value: string) => {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? null : date;
};

const getSafeClicks = (value: number) => (Number.isFinite(value) ? Math.max(0, value) : 0);

const getPointX = (index: number, pointsCount: number) => {
  const start = PADDING.left;
  const end = CHART_WIDTH - PADDING.right;

  if (pointsCount <= 1) {
    return start + (end - start) / 2;
  }

  return start + (index / (pointsCount - 1)) * (end - start);
};

const getPointY = (clicks: number, maxClicks: number) => {
  const start = PADDING.top;
  const end = CHART_HEIGHT - PADDING.bottom;
  const ratio = maxClicks > 0 ? clicks / maxClicks : 0;

  return end - ratio * (end - start);
};

const formatNumber = (value: number) => numberFormatter.format(value);

const formatDateLabel = (value: string, groupBy: AnalyticsGroupBy) => {
  const date = parseDate(value);

  if (!date) {
    return "—";
  }

  return groupBy === "hour" ? hourFormatter.format(date) : dayFormatter.format(date);
};

const formatTooltipDate = (value: string) => {
  const date = parseDate(value);
  return date ? fullDateFormatter.format(date) : "Дата не определена";
};

const getVisibleLabelIndexes = (pointsCount: number) => {
  if (pointsCount <= MAX_X_LABELS) {
    return Array.from({ length: pointsCount }, (_, index) => index);
  }

  const labelsCount = MAX_X_LABELS;
  const indexes = new Set<number>();

  for (let index = 0; index < labelsCount; index += 1) {
    indexes.add(Math.round((index * (pointsCount - 1)) / (labelsCount - 1)));
  }

  return Array.from(indexes).sort((a, b) => a - b);
};

const buildLinePath = (points: ClicksChartPoint[]) =>
  points.map((point, index) => `${index === 0 ? "M" : "L"} ${point.x} ${point.y}`).join(" ");

const buildAreaPath = (points: ClicksChartPoint[]) => {
  if (points.length < 2) {
    return "";
  }

  const baseline = CHART_HEIGHT - PADDING.bottom;
  const linePath = buildLinePath(points);
  const firstPoint = points[0];
  const lastPoint = points[points.length - 1];

  return `${linePath} L ${lastPoint.x} ${baseline} L ${firstPoint.x} ${baseline} Z`;
};

const buildYAxisLabels = (maxClicks: number): ChartAxisLabel[] => {
  const values =
    maxClicks <= 1
      ? [maxClicks, 0]
      : [maxClicks, Math.round(maxClicks / 2), 0].filter((value, index, array) => array.indexOf(value) === index);

  return values.map((value) => ({
    key: `y-${value}`,
    label: formatNumber(value),
    y: getPointY(value, maxClicks),
  }));
};

export const buildClicksChartData = (
  series: TimeSeriesPoint[],
  groupBy: AnalyticsGroupBy,
): ClicksChartData => {
  const pointsCount = series.length;
  const maxClicks = Math.max(1, ...series.map((point) => getSafeClicks(point.clicks)));
  const points = series.map((point, index) => {
    const clicks = getSafeClicks(point.clicks);
    const xLabel = formatDateLabel(point.periodStart, groupBy);

    return {
      key: `${point.periodStart || "empty"}-${index}`,
      x: getPointX(index, pointsCount),
      y: getPointY(clicks, maxClicks),
      clicks,
      xLabel,
      tooltip: `${formatTooltipDate(point.periodStart)}: ${formatNumber(clicks)}`,
    };
  });

  return {
    width: CHART_WIDTH,
    height: CHART_HEIGHT,
    linePath: buildLinePath(points),
    areaPath: buildAreaPath(points),
    points,
    xAxisLabels: getVisibleLabelIndexes(points.length).map((index) => ({
      key: `x-${index}`,
      label: points[index]?.xLabel ?? "—",
      x: points[index]?.x,
    })),
    yAxisLabels: buildYAxisLabels(maxClicks),
  };
};
