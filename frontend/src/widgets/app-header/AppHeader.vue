<script setup lang="ts">
import { computed } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { UiBrandLogo, UiButton } from "@/shared/ui";
import { ROUTES } from "@/shared/lib/routes/paths";

withDefaults(
  defineProps<{
    userEmail?: string;
    loading?: boolean;
  }>(),
  {
    userEmail: "",
    loading: false,
  },
);

const emit = defineEmits<{
  logout: [];
}>();

const route = useRoute();

const isDashboardActive = computed(
  () => route.path.startsWith(ROUTES.dashboard) || route.path.startsWith("/links"),
);

const onLogout = () => emit("logout");
</script>

<template>
  <header class="app-header">
    <div class="app-header__left">
      <RouterLink :to="ROUTES.dashboard" class="app-header__brand-link" aria-label="TrackLink dashboard">
        <UiBrandLogo variant="header" />
      </RouterLink>

      <RouterLink
        :to="ROUTES.dashboard"
        class="app-header__dashboard-link"
        :class="{ 'is-active': isDashboardActive }"
      >
        Dashboard
      </RouterLink>
    </div>

    <div class="app-header__right">
      <span v-if="loading" class="app-header__email app-header__email--muted">Loading...</span>
      <span v-else-if="userEmail" class="app-header__email">{{ userEmail }}</span>
      <span v-else class="app-header__email app-header__email--muted">No email</span>

      <RouterLink :to="ROUTES.settings" class="app-header__settings-link">Settings</RouterLink>

      <UiButton variant="ghost" size="sm" @click="onLogout">
        <span>Logout</span>
        <template #iconRight>
          <svg viewBox="0 0 24 24" class="app-header__logout-icon" aria-hidden="true">
            <path
              d="M16 17v-1.5h-6v-7h6V7l4 5-4 5Zm-1.5 3H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h8.5v1.5H6a.5.5 0 0 0-.5.5v12c0 .3.2.5.5.5h8.5V20Z"
              fill="currentColor"
            />
          </svg>
        </template>
      </UiButton>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  min-height: 72px;
  background: var(--tl-color-white);
  border-bottom: 1px solid rgb(37 31 63 / 12%);
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.app-header__left,
.app-header__right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-header__brand-link {
  display: inline-flex;
  text-decoration: none;
}

.app-header__dashboard-link {
  text-decoration: none;
  color: var(--tl-color-text);
  font-weight: 600;
  font-size: 14px;
  padding: 10px 14px;
  border-radius: 999px;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.app-header__dashboard-link:hover {
  background: var(--tl-color-surface-muted);
}

.app-header__dashboard-link.is-active {
  background: var(--tl-color-primary);
  color: var(--tl-color-white);
}

.app-header__email {
  font-size: 14px;
  color: var(--tl-color-text);
  white-space: nowrap;
}

.app-header__email--muted {
  color: var(--tl-color-text-muted);
}

.app-header__settings-link {
  text-decoration: none;
  color: var(--tl-color-primary);
  font-weight: 600;
  font-size: 14px;
}

.app-header__logout-icon {
  width: 14px;
  height: 14px;
}

@media (max-width: 1023px) {
  .app-header {
    padding-inline: 24px;
  }
}

@media (max-width: 767px) {
  .app-header {
    padding-inline: 16px;
    align-items: flex-start;
  }

  .app-header__right {
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .app-header__email {
    display: none;
  }
}
</style>
