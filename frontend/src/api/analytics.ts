import { http } from "@/api/http";
import type { Link } from "@/entities/link/link.types";

export interface TimeSeriesPoint {
  periodStart: string;
  clicks: number;
}

export interface LinkAnalyticsResponse {
  linkId: string;
  totalClicks: number;
  clicksLast24h: number;
  lastClickedAt?: string | null;
  series: TimeSeriesPoint[];
}

export interface ClickEvent {
  id: string;
  linkId: string;
  clickedAt: string;
  referrer?: string | null;
  userAgent?: string | null;
}

export interface RecentClicksResponse {
  items: ClickEvent[];
}

export interface GetLinkAnalyticsParams {
  from?: string;
  to?: string;
  groupBy?: "hour" | "day";
}

export interface ListRecentClicksParams {
  limit?: number;
}

export interface DashboardResponse {
  totalLinks: number;
  activeLinks: number;
  totalClicks: number;
  clicksLast24h: number;
  recentLinks: Link[];
}

export const getDashboard = () => http.get<DashboardResponse>("/dashboard");

export const getLinkAnalytics = (linkId: string, params: GetLinkAnalyticsParams = {}) => {
  const searchParams = new URLSearchParams();

  if (params.from) {
    searchParams.set("from", params.from);
  }

  if (params.to) {
    searchParams.set("to", params.to);
  }

  if (params.groupBy) {
    searchParams.set("groupBy", params.groupBy);
  }

  const queryString = searchParams.toString();
  return http.get<LinkAnalyticsResponse>(
    `/links/${encodeURIComponent(linkId)}/analytics${queryString ? `?${queryString}` : ""}`,
  );
};

export const listRecentClicks = (linkId: string, params: ListRecentClicksParams = {}) => {
  const searchParams = new URLSearchParams();

  if (params.limit !== undefined) {
    searchParams.set("limit", String(params.limit));
  }

  const queryString = searchParams.toString();
  return http.get<RecentClicksResponse>(
    `/links/${encodeURIComponent(linkId)}/clicks${queryString ? `?${queryString}` : ""}`,
  );
};
