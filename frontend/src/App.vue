<script setup lang="ts">
import { computed } from "vue";
import { RouterView, useRoute } from "vue-router";
import { useLogout } from "@/features/auth/logout/useLogout";
import AuthLayout from "@/layouts/AuthLayout.vue";
import AppLayout from "@/layouts/AppLayout.vue";
import { useSession } from "@/entities/session/useSession";
import { useToast } from "@/shared/composables/useToast";
import { useI18n } from "@/shared/composables/useI18n";
import { UiPageState, UiToast } from "@/shared/ui";

const route = useRoute();
const session = useSession();
const toast = useToast();
const { t } = useI18n();
const { isLoggingOut, logout } = useLogout();

const layoutComponent = computed(() => (route.meta.layout === "auth" ? AuthLayout : AppLayout));
const isSessionLoading = computed(() => session.isSessionLoading.value);

const handleLogout = () => logout();
</script>

<template>
  <UiPageState
    v-if="isSessionLoading"
    class="app-session-state"
    type="loading"
    :title="t('session.checking.title')"
    :description="t('session.checking.description')"
  />
  <component :is="layoutComponent" v-else :user-email="session.user.value?.email" :loading="isLoggingOut" @logout="handleLogout">
    <RouterView />
  </component>
  <UiToast :items="toast.toasts.value" @remove="toast.remove" />
</template>

<style scoped>
.app-session-state {
  min-height: 100dvh;
  justify-content: center;
}
</style>
