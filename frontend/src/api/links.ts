import { http } from "@/api/http";
import type { LinkListResponse, ListLinksParams } from "@/entities/link/link.types";

const appendNumberParam = (searchParams: URLSearchParams, key: string, value?: number) => {
  if (value !== undefined) {
    searchParams.set(key, String(value));
  }
};

export const listLinks = (params: ListLinksParams = {}) => {
  const searchParams = new URLSearchParams();
  const query = params.q?.trim();

  appendNumberParam(searchParams, "page", params.page);
  appendNumberParam(searchParams, "pageSize", params.pageSize);

  if (query) {
    searchParams.set("q", query);
  }

  if (params.status) {
    searchParams.set("status", params.status);
  }

  const queryString = searchParams.toString();
  return http.get<LinkListResponse>(`/links${queryString ? `?${queryString}` : ""}`);
};
