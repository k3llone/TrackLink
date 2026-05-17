<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiPageHeader, UiPageState } from "@/shared/ui";

const route = useRoute();

const linkId = computed(() => {
  const rawId = route.params.id;
  return Array.isArray(rawId) ? rawId[0] : rawId;
});
</script>

<template>
  <section class="admin-link-analytics-placeholder-page">
    <UiPageHeader
      title="Аналитика ссылки"
      :subtitle="linkId ? `Admin analytics entry point для ссылки ${linkId}.` : 'Ссылка не выбрана.'"
      :back-to="ROUTES.admin"
      back-label="Admin panel"
    />

    <UiPageState
      v-if="linkId"
      type="empty"
      title="Аналитика пока недоступна"
      description="Для полноценной admin analytics page нужен отдельный backend endpoint. Эта страница оставлена как admin entry point."
    />

    <UiPageState
      v-else
      type="not-found"
      title="Ссылка не выбрана"
      description="Откройте ссылку из административной таблицы."
      :action-to="ROUTES.admin"
      action-text="Вернуться в admin panel"
    />
  </section>
</template>

<style scoped>
.admin-link-analytics-placeholder-page {
  width: 100%;
}
</style>
