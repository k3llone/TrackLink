import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";
import RouteStubPage from "@/pages/RouteStubPage.vue";

const routes: RouteRecordRaw[] = [
  { path: "/login", name: "login", component: RouteStubPage, props: { title: "Login" } },
  { path: "/register", name: "register", component: RouteStubPage, props: { title: "Register" } },
  { path: "/dashboard", name: "dashboard", component: RouteStubPage, props: { title: "Dashboard" } },
  { path: "/links/new", name: "link-new", component: RouteStubPage, props: { title: "Create Link" } },
  { path: "/links/:id", name: "link-details", component: RouteStubPage, props: { title: "Link Details" } },
  {
    path: "/links/:id/edit",
    name: "link-edit",
    component: RouteStubPage,
    props: { title: "Edit Link" },
  },
  { path: "/profile", name: "profile", component: RouteStubPage, props: { title: "Profile" } },
  { path: "/admin", name: "admin", component: RouteStubPage, props: { title: "Admin" } },
  {
    path: "/:code/error",
    name: "code-error",
    component: RouteStubPage,
    props: (route) => ({ title: `Redirect Error: ${String(route.params.code ?? "")}` }),
  },
  { path: "/error", name: "error", component: RouteStubPage, props: { title: "Application Error" } },
  { path: "/:pathMatch(.*)*", redirect: "/error" },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
