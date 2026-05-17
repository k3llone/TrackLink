import { http } from "@/api/http";
import type { Link, Pagination } from "@/entities/link/link.types";

export type AdminLink = Link & {
  ownerEmail?: string | null;
};

export interface AdminLinkListResponse {
  items: AdminLink[];
  pagination: Pagination;
}

export interface ListAdminLinksParams {
  page?: number;
  pageSize?: number;
  q?: string;
}

const appendNumberParam = (searchParams: URLSearchParams, key: string, value?: number) => {
  if (value !== undefined) {
    searchParams.set(key, String(value));
  }
};

export const listAdminLinks = (params: ListAdminLinksParams = {}) => {
  const searchParams = new URLSearchParams();
  const query = params.q?.trim();

  appendNumberParam(searchParams, "page", params.page);
  appendNumberParam(searchParams, "pageSize", params.pageSize);

  if (query) {
    searchParams.set("q", query);
  }

  const queryString = searchParams.toString();
  return http.get<AdminLinkListResponse>(`/admin/links${queryString ? `?${queryString}` : ""}`);
};

export const blockAdminLink = (linkId: string) =>
  http.patch<AdminLink>(`/admin/links/${encodeURIComponent(linkId)}/block`);

export const unblockAdminLink = (linkId: string) =>
  http.patch<AdminLink>(`/admin/links/${encodeURIComponent(linkId)}/unblock`);
