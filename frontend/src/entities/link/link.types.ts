export type LinkStatus = "active" | "inactive" | "blocked" | "deleted";

export interface Link {
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

export interface Pagination {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}

export interface LinkListResponse {
  items: Link[];
  pagination: Pagination;
}

export interface CreateLinkRequest {
  targetUrl: string;
  customAlias?: string | null;
}

export interface ListLinksParams {
  page?: number;
  pageSize?: number;
  q?: string;
}
