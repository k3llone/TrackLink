<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { getDashboard, type DashboardResponse } from "@/api/analytics";
import type { ApiClientError } from "@/api/types";
import { DashboardSummary, RecentLinks } from "@/widgets/dashboard";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiButton, UiPageHeader, UiPageState } from "@/shared/ui";

const router = useRouter();

const dashboard = ref<DashboardResponse | null>(null);
const isLoading = ref(false);
const errorMessage = ref("");

const emptySummary: DashboardResponse = {
  totalLinks: 0,
  activeLinks: 0,
  totalClicks: 0,
  clicksLast24h: 0,
  recentLinks: [],
};

const summary = computed(() => dashboard.value ?? emptySummary);
const isEmpty = computed(() => !isLoading.value && !errorMessage.value && dashboard.value?.totalLinks === 0);

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const getDashboardErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 401) {
      return "Сессия недействительна. Войдите заново, чтобы открыть dashboard.";
    }

    if (error.status === 403) {
      return "У вас нет доступа к dashboard этого аккаунта.";
    }
  }

  return "Не удалось загрузить dashboard. Проверьте соединение и повторите попытку.";
};

const loadDashboard = async () => {
  if (isLoading.value) {
    return;
  }

  isLoading.value = true;
  errorMessage.value = "";

  try {
    dashboard.value = await getDashboard();
  } catch (error: unknown) {
    dashboard.value = null;
    errorMessage.value = getDashboardErrorMessage(error);
  } finally {
    isLoading.value = false;
  }
};

const goToCreateLink = () => {
  void router.push(ROUTES.linkCreate);
};

onMounted(() => {
  void loadDashboard();
});
</script>

<template>
  <section class="dashboard-page">
    <UiPageHeader title="Dashboard" subtitle="Ключевые показатели аккаунта и последние созданные ссылки." />

    <DashboardSummary
      v-if="isLoading || dashboard"
      :total-links="summary.totalLinks"
      :active-links="summary.activeLinks"
      :total-clicks="summary.totalClicks"
      :clicks-last24h="summary.clicksLast24h"
      :loading="isLoading"
    />

    <UiPageState
      v-if="isLoading && !dashboard"
      type="loading"
      title="Загружаем dashboard"
      description="Получаем метрики и последние ссылки аккаунта."
    />

    <UiPageState
      v-else-if="errorMessage"
      type="error"
      title="Dashboard недоступен"
      :description="errorMessage"
      action-text="Повторить"
      @action="loadDashboard"
    />

    <UiPageState
      v-else-if="isEmpty"
      type="empty"
      title="Ссылок пока нет"
      description="Создайте первую короткую ссылку, чтобы здесь появились метрики и последние действия."
      action-text="Создать первую ссылку"
      @action="goToCreateLink"
    />

    <div v-else-if="dashboard" class="dashboard-page__content">
      <section class="dashboard-page__quick-actions" aria-labelledby="dashboard-quick-actions-title">
        <div class="dashboard-page__quick-copy">
          <h2 id="dashboard-quick-actions-title" class="dashboard-page__section-title">Быстрые действия</h2>
          <p class="dashboard-page__section-description">
            Создайте новую короткую ссылку и сразу начните собирать статистику переходов.
          </p>
        </div>

        <UiButton type="button" @click="goToCreateLink">Создать ссылку</UiButton>
      </section>

      <RecentLinks :links="dashboard.recentLinks" />
    </div>
  </section>
</template>

<style scoped>
.dashboard-page {
  width: 100%;
}

.dashboard-page__content {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.dashboard-page__quick-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 22px 0;
  border-top: 1px solid rgb(37 31 63 / 10%);
  border-bottom: 1px solid rgb(37 31 63 / 10%);
}

.dashboard-page__quick-copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.dashboard-page__section-title {
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.dashboard-page__section-description {
  color: var(--tl-color-text-muted);
  font-size: 14px;
  max-width: 560px;
}

@media (max-width: 767px) {
  .dashboard-page__content {
    gap: 24px;
  }

  .dashboard-page__quick-actions {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
