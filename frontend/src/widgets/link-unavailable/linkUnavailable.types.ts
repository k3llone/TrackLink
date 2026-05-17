import type { LinkUnavailableRouteReason } from "@/shared/lib/routes/paths";

export type LinkUnavailableReason = LinkUnavailableRouteReason;

export const LINK_UNAVAILABLE_REASONS: readonly LinkUnavailableReason[] = [
  "not_found",
  "inactive",
  "blocked",
  "deleted",
  "gone",
];

export const isLinkUnavailableReason = (value: string): value is LinkUnavailableReason =>
  LINK_UNAVAILABLE_REASONS.includes(value as LinkUnavailableReason);

export const LINK_UNAVAILABLE_REASON_BY_SLUG: Readonly<Record<string, LinkUnavailableReason>> = {
  "not-found": "not_found",
  inactive: "inactive",
  blocked: "blocked",
  deleted: "deleted",
  gone: "gone",
};

export const normalizeLinkUnavailableReason = (value: unknown): LinkUnavailableReason => {
  if (typeof value !== "string") {
    return "not_found";
  }

  if (isLinkUnavailableReason(value)) {
    return value;
  }

  return LINK_UNAVAILABLE_REASON_BY_SLUG[value] ?? "not_found";
};
