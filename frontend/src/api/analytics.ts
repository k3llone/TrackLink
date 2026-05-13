import { http } from "@/api/http";

export type LinkStatus = "active" | "inactive" | "blocked" | "deleted";

export interface DashboardLink {
  id: string;
  ownerId: string;
  code: string;
  customAlias?: string | null;
  shortUrl: string;
  targetUrl: string;
  status: LinkStatus;
  totalClicks: number;
  lastClickedAt?: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface DashboardResponse {
  totalLinks: number;
  activeLinks: number;
  totalClicks: number;
  clicksLast24h: number;
  recentLinks: DashboardLink[];
}

export const getDashboard = () => http.get<DashboardResponse>("/dashboard");
