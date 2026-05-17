import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";
import LoginPage from "@/pages/auth/login/LoginPage.vue";
import RegisterPage from "@/pages/auth/register/RegisterPage.vue";
import ForgotPasswordPage from "@/pages/auth/forgot-password/ForgotPasswordPage.vue";
import ResetPasswordPage from "@/pages/auth/reset-password/ResetPasswordPage.vue";
import TwoFactorPage from "@/pages/auth/two-factor/TwoFactorPage.vue";
import DashboardPage from "@/pages/dashboard/DashboardPage.vue";
import CreateLinkPage from "@/pages/links/create/CreateLinkPage.vue";
import LinkAnalyticsPage from "@/pages/links/analytics/LinkAnalyticsPage.vue";
import EditLinkPage from "@/pages/links/edit/EditLinkPage.vue";
import SettingsPage from "@/pages/settings/SettingsPage.vue";
import AdminPage from "@/pages/admin/AdminPage.vue";
import AdminLinkAnalyticsPlaceholderPage from "@/pages/admin/AdminLinkAnalyticsPlaceholderPage.vue";
import LinkUnavailablePage from "@/pages/public/LinkUnavailablePage.vue";
import NotFoundPage from "@/pages/not-found/NotFoundPage.vue";
import { useSession } from "@/entities/session/useSession";
import { ROUTES } from "@/shared/lib/routes/paths";

type AppRouteMeta = {
  layout: "auth" | "app";
  requiresAuth?: boolean;
  requiresGuest?: boolean;
  requiresAdmin?: boolean;
};

declare module "vue-router" {
  interface RouteMeta {
    layout: AppRouteMeta["layout"];
    requiresAuth?: AppRouteMeta["requiresAuth"];
    requiresGuest?: AppRouteMeta["requiresGuest"];
    requiresAdmin?: AppRouteMeta["requiresAdmin"];
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    redirect: ROUTES.dashboard,
  },
  {
    path: ROUTES.login,
    name: "login",
    component: LoginPage,
    meta: { layout: "auth", requiresGuest: true },
  },
  {
    path: ROUTES.register,
    name: "register",
    component: RegisterPage,
    meta: { layout: "auth", requiresGuest: true },
  },
  {
    path: ROUTES.forgotPassword,
    name: "forgot-password",
    component: ForgotPasswordPage,
    meta: { layout: "auth", requiresGuest: true },
  },
  {
    path: ROUTES.resetPassword,
    name: "reset-password",
    component: ResetPasswordPage,
    meta: { layout: "auth", requiresGuest: true },
  },
  {
    path: ROUTES.twoFactor,
    name: "two-factor",
    component: TwoFactorPage,
    meta: { layout: "auth", requiresGuest: true },
  },
  {
    path: ROUTES.linkUnavailableStatus,
    name: "link-unavailable-status",
    component: LinkUnavailablePage,
    meta: { layout: "auth" },
  },
  {
    path: ROUTES.linkUnavailable,
    name: "link-unavailable",
    component: LinkUnavailablePage,
    meta: { layout: "auth" },
  },
  {
    path: ROUTES.dashboard,
    name: "dashboard",
    component: DashboardPage,
    meta: { layout: "app", requiresAuth: true },
  },
  {
    path: ROUTES.linkCreate,
    name: "link-create",
    component: CreateLinkPage,
    meta: { layout: "app", requiresAuth: true },
  },
  {
    path: ROUTES.linkDetails,
    name: "link-details",
    component: LinkAnalyticsPage,
    meta: { layout: "app", requiresAuth: true },
  },
  {
    path: ROUTES.linkEdit,
    name: "link-edit",
    component: EditLinkPage,
    meta: { layout: "app", requiresAuth: true },
  },
  {
    path: ROUTES.settings,
    name: "settings",
    component: SettingsPage,
    meta: { layout: "app", requiresAuth: true },
  },
  {
    path: ROUTES.admin,
    name: "admin",
    component: AdminPage,
    meta: { layout: "app", requiresAuth: true, requiresAdmin: true },
  },
  {
    path: ROUTES.adminLinkDetails,
    name: "admin-link-details",
    component: AdminLinkAnalyticsPlaceholderPage,
    meta: { layout: "app", requiresAuth: true, requiresAdmin: true },
  },
  {
    path: ROUTES.notFound,
    name: "not-found",
    component: NotFoundPage,
    meta: { layout: "auth" },
  },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to) => {
  const session = useSession();
  const shouldResolveSession = to.meta.requiresAuth || to.meta.requiresGuest || to.meta.requiresAdmin;

  if (shouldResolveSession) {
    await session.loadCurrentUser();
  }

  if (to.meta.requiresAuth && !session.isAuthenticated.value) {
    return { path: ROUTES.login };
  }

  if (to.meta.requiresGuest && session.isAuthenticated.value) {
    return { path: ROUTES.dashboard };
  }

  if (to.meta.requiresAdmin && !session.isAdmin.value) {
    return { path: ROUTES.dashboard };
  }

  return true;
});
