<script setup lang="ts">
import { computed } from "vue";
import { RouterView } from "vue-router";
import { useRoute } from "vue-router";
import AuthLayout from "@/layouts/AuthLayout.vue";
import AppLayout from "@/layouts/AppLayout.vue";
import { ROUTES } from "@/shared/lib/routes/paths";

const route = useRoute();

const layoutComponent = computed(() => (route.meta.layout === "auth" ? AuthLayout : AppLayout));

const handleLogout = () => {
  // TODO: replace with real logout flow from session/auth feature.
  window.location.assign(ROUTES.login);
};
</script>

<template>
  <component :is="layoutComponent" @logout="handleLogout">
    <RouterView />
  </component>
</template>
