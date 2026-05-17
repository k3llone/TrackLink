<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import { useI18n } from "@/shared/composables/useI18n";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiPageHeader, UiPageState } from "@/shared/ui";

const route = useRoute();
const { t } = useI18n();

const linkId = computed(() => {
  const rawId = route.params.id;
  return Array.isArray(rawId) ? rawId[0] : rawId;
});

const subtitle = computed(() =>
  linkId.value ? t("admin.placeholder.subtitleWithId", { linkId: linkId.value }) : t("admin.placeholder.subtitleNoLink"),
);
</script>

<template>
  <section class="admin-link-analytics-placeholder-page">
    <UiPageHeader
      :title="t('analytics.page.title')"
      :subtitle="subtitle"
      :back-to="ROUTES.admin"
      :back-label="t('common.adminPanel')"
    />

    <UiPageState
      v-if="linkId"
      type="empty"
      :title="t('admin.placeholder.emptyTitle')"
      :description="t('admin.placeholder.emptyDescription')"
    />

    <UiPageState
      v-else
      type="not-found"
      :title="t('admin.placeholder.noLinkTitle')"
      :description="t('admin.placeholder.noLinkDescription')"
      :action-to="ROUTES.admin"
      :action-text="t('admin.placeholder.noLinkAction')"
    />
  </section>
</template>

<style scoped>
.admin-link-analytics-placeholder-page {
  width: 100%;
}
</style>
