import { http } from "@/api/http";
import type { Link } from "@/entities/link/link.types";

export interface DashboardResponse {
  totalLinks: number;
  activeLinks: number;
  totalClicks: number;
  clicksLast24h: number;
  recentLinks: Link[];
}

export const getDashboard = () => http.get<DashboardResponse>("/dashboard");
