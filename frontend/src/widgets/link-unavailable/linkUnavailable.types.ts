export type LinkUnavailableReason = "not_found" | "inactive" | "blocked" | "deleted" | "gone";

export const LINK_UNAVAILABLE_REASONS: readonly LinkUnavailableReason[] = [
  "not_found",
  "inactive",
  "blocked",
  "deleted",
  "gone",
];

export const isLinkUnavailableReason = (value: string): value is LinkUnavailableReason =>
  LINK_UNAVAILABLE_REASONS.includes(value as LinkUnavailableReason);
