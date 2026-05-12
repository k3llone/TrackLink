<script setup lang="ts">
import { computed } from "vue";
import { RouterView, useRoute } from "vue-router";
import { useLogout } from "@/features/auth/logout/useLogout";
import AuthLayout from "@/layouts/AuthLayout.vue";
import AppLayout from "@/layouts/AppLayout.vue";
import { useSession } from "@/entities/session/useSession";
import { useToast } from "@/shared/composables/useToast";
import { UiToast } from "@/shared/ui";

const route = useRoute();
const session = useSession();
const toast = useToast();
const { isLoggingOut, logout } = useLogout();

const layoutComponent = computed(() => (route.meta.layout === "auth" ? AuthLayout : AppLayout));

const handleLogout = () => logout();
</script>

<template>
  <component :is="layoutComponent" :user-email="session.user.value?.email" :loading="isLoggingOut" @logout="handleLogout">
    <RouterView />
  </component>
  <UiToast :items="toast.toasts.value" @remove="toast.remove" />
</template>
