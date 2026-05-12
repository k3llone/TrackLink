<script setup lang="ts">
import { computed } from "vue";
import { RouterView, useRoute } from "vue-router";
import AuthLayout from "@/layouts/AuthLayout.vue";
import AppLayout from "@/layouts/AppLayout.vue";
import { useSession } from "@/entities/session/useSession";
import { useToast } from "@/shared/composables/useToast";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiToast } from "@/shared/ui";

const route = useRoute();
const session = useSession();
const toast = useToast();

const layoutComponent = computed(() => (route.meta.layout === "auth" ? AuthLayout : AppLayout));

const handleLogout = () => {
  // TODO: replace with real logout flow from session/auth feature.
  session.clearSession();
  window.location.assign(ROUTES.login);
};
</script>

<template>
  <component :is="layoutComponent" :user-email="session.user.value?.email" @logout="handleLogout">
    <RouterView />
  </component>
  <UiToast :items="toast.toasts.value" @remove="toast.remove" />
</template>
