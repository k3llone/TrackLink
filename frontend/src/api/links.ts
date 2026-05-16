import { http } from "@/api/http";
import type {
  CreateLinkRequest,
  Link,
  LinkListResponse,
  ListLinksParams,
  UpdateLinkStatusRequest,
} from "@/entities/link/link.types";

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

export const findOwnLinkById = async (linkId: string) => {
  const normalizedLinkId = linkId.trim();

  if (!normalizedLinkId) {
    return null;
  }

  const pageSize = 100;
  let page = 1;
  let hasMorePages = true;

  while (hasMorePages) {
    const response = await listLinks({ page, pageSize });
    const link = response.items.find((item) => item.id === normalizedLinkId);

    if (link) {
      return link;
    }

    hasMorePages = Boolean(response.pagination.totalPages && page < response.pagination.totalPages);
    page += 1;
  }

  return null;
};

export const createLink = (payload: CreateLinkRequest) => http.post<Link>("/links", payload);

export const updateLinkStatus = (linkId: string, payload: UpdateLinkStatusRequest) =>
  http.patch<Link>(`/links/${encodeURIComponent(linkId)}/status`, payload);

export const deleteLink = (linkId: string) => http.delete<void>(`/links/${encodeURIComponent(linkId)}`);
