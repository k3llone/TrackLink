export const ROUTES = {
  login: "/login",
  register: "/register",
  forgotPassword: "/forgot-password",
  resetPassword: "/reset-password",
  twoFactor: "/2fa",

  dashboard: "/dashboard",
  linkCreate: "/links/create",
  linkDetails: "/links/:id",
  linkEdit: "/links/:id/edit",
  settings: "/settings",
  admin: "/admin",

  notFound: "/:pathMatch(.*)*",
} as const;

export const getLinkDetailsPath = (id: string) => `/links/${id}`;
export const getLinkEditPath = (id: string) => `/links/${id}/edit`;
