<script setup lang="ts">
import { computed } from "vue";
import { AppHeader } from "@/widgets/app-header";

const props = withDefaults(
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

const resolvedUserEmail = computed(() => props.userEmail);

const onLogout = () => emit("logout");
</script>

<template>
  <div class="app-layout">
    <AppHeader :user-email="resolvedUserEmail" :loading="loading" @logout="onLogout" />

    <main class="app-layout__main">
      <div class="app-layout__container">
        <slot />
      </div>
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100dvh;
  background: var(--tl-color-bg);
}

.app-layout__main {
  padding: 24px 16px 40px;
}

.app-layout__container {
  width: 100%;
  max-width: var(--tl-page-max-width);
  margin: 0 auto;
}

@media (max-width: 1023px) {
  .app-layout__main {
    padding-inline: 24px;
  }
}

@media (max-width: 767px) {
  .app-layout__main {
    padding-inline: 16px;
  }
}
</style>
