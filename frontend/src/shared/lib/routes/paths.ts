export const ROUTES = {
  login: "/login",
  register: "/register",
  forgotPassword: "/forgot-password",
  resetPassword: "/reset-password",
  twoFactor: "/2fa",
  linkUnavailable: "/link-unavailable",
  linkUnavailableStatus: "/link-unavailable/:reason",

  dashboard: "/dashboard",
  linkCreate: "/links/create",
  linkDetails: "/links/:id",
  linkEdit: "/links/:id/edit",
  settings: "/settings",
  admin: "/admin",
  adminLinkDetails: "/admin/links/:id",

  notFound: "/:pathMatch(.*)*",
} as const;

export type LinkUnavailableRouteReason = "not_found" | "inactive" | "blocked" | "deleted" | "gone";

const linkUnavailableReasonSlugs: Record<LinkUnavailableRouteReason, string> = {
  not_found: "not-found",
  inactive: "inactive",
  blocked: "blocked",
  deleted: "deleted",
  gone: "gone",
};

export const getLinkDetailsPath = (id: string) => `/links/${id}`;
export const getLinkEditPath = (id: string) => `/links/${id}/edit`;
export const getAdminLinkDetailsPath = (id: string) => `/admin/links/${encodeURIComponent(id)}`;
export const getLinkUnavailablePath = (reason: LinkUnavailableRouteReason) =>
  `${ROUTES.linkUnavailable}/${linkUnavailableReasonSlugs[reason]}`;
