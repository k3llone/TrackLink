<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getLinkDetailsPath, ROUTES } from "@/shared/lib/routes/paths";
import { UiButton, UiPageHeader, UiPageState, UiStatusBadge } from "@/shared/ui";

const route = useRoute();
const router = useRouter();

const linkId = computed(() => {
  const rawId = route.params.id;
  return Array.isArray(rawId) ? rawId[0] : rawId || "";
});

const analyticsPath = computed(() => (linkId.value ? getLinkDetailsPath(linkId.value) : ROUTES.dashboard));

const goToAnalytics = () => {
  if (linkId.value) {
    void router.push(analyticsPath.value);
  }
};

const goToDashboard = () => {
  void router.push(ROUTES.dashboard);
};
</script>

<template>
  <section class="edit-link-page">
    <UiPageHeader
      title="Редактирование ссылки"
      subtitle="Изменение target URL запланировано для Post-MVP и недоступно в текущей версии."
      :back-to="ROUTES.dashboard"
      back-label="Dashboard"
    >
      <UiStatusBadge status="pending" label="Post-MVP" />
    </UiPageHeader>

    <UiPageState
      v-if="!linkId"
      type="not-found"
      title="Ссылка не выбрана"
      description="Откройте нужную ссылку из dashboard или списка ссылок."
      :action-to="ROUTES.dashboard"
      action-text="Вернуться на dashboard"
    />

    <div v-else class="edit-link-page__content">
      <section class="edit-link-page__notice" aria-labelledby="edit-link-notice-title">
        <div class="edit-link-page__notice-copy">
          <h2 id="edit-link-notice-title" class="edit-link-page__notice-title">Target URL нельзя изменить в MVP</h2>
          <p class="edit-link-page__notice-text">
            Изменение назначения уже распространенной короткой ссылки несет риск подмены целевого ресурса. Форма
            редактирования будет подключена после появления backend/API-поддержки, истории изменений и дополнительных
            ограничений безопасности.
          </p>
        </div>

        <dl class="edit-link-page__meta" aria-label="Данные выбранной ссылки">
          <div class="edit-link-page__meta-row">
            <dt>Link ID</dt>
            <dd>{{ linkId }}</dd>
          </div>
          <div class="edit-link-page__meta-row">
            <dt>Target URL edit</dt>
            <dd>Post-MVP</dd>
          </div>
        </dl>
      </section>

      <div class="edit-link-page__actions">
        <UiButton type="button" @click="goToAnalytics">К аналитике</UiButton>
        <UiButton type="button" variant="secondary" @click="goToDashboard">Dashboard</UiButton>
      </div>
    </div>
  </section>
</template>

<style scoped>
.edit-link-page {
  width: 100%;
}

.edit-link-page__content {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.edit-link-page__notice {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(220px, 320px);
  gap: 20px;
  padding: 22px 0;
  border-top: 1px solid rgb(37 31 63 / 10%);
  border-bottom: 1px solid rgb(37 31 63 / 10%);
}

.edit-link-page__notice-copy {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.edit-link-page__notice-title {
  margin: 0;
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.edit-link-page__notice-text {
  margin: 0;
  max-width: 720px;
  color: var(--tl-color-text-muted);
  font-size: 14px;
  line-height: 1.6;
}

.edit-link-page__meta {
  display: grid;
  gap: 12px;
  margin: 0;
  padding: 0;
}

.edit-link-page__meta-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.edit-link-page__meta-row dt {
  color: var(--tl-color-text-muted);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
}

.edit-link-page__meta-row dd {
  margin: 0;
  color: var(--tl-color-text);
  font-size: 14px;
  font-weight: 700;
  overflow-wrap: anywhere;
}

.edit-link-page__actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

@media (max-width: 767px) {
  .edit-link-page__notice {
    grid-template-columns: 1fr;
  }

  .edit-link-page__actions,
  .edit-link-page__actions :deep(.ui-button) {
    width: 100%;
  }
}
</style>
