<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { getDashboard, type DashboardResponse } from "@/api/analytics";
import { listLinks } from "@/api/links";
import type { ApiClientError } from "@/api/types";
import type { Link, Pagination } from "@/entities/link/link.types";
import { DashboardSummary } from "@/widgets/dashboard";
import { LinksTable } from "@/widgets/links-table";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiButton, UiPageHeader, UiPageState } from "@/shared/ui";

const router = useRouter();

const dashboard = ref<DashboardResponse | null>(null);
const isDashboardLoading = ref(false);
const dashboardErrorMessage = ref("");
const links = ref<Link[]>([]);
const linksPagination = ref<Pagination | null>(null);
const isLinksLoading = ref(false);
const linksErrorMessage = ref("");
const linksPage = ref(1);
const linksPageSize = 20;
const linksQ = ref("");
let linksRequestId = 0;

const emptySummary: DashboardResponse = {
  totalLinks: 0,
  activeLinks: 0,
  totalClicks: 0,
  clicksLast24h: 0,
  recentLinks: [],
};

const summary = computed(() => dashboard.value ?? emptySummary);

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
  if (isDashboardLoading.value) {
    return;
  }

  isDashboardLoading.value = true;
  dashboardErrorMessage.value = "";

  try {
    dashboard.value = await getDashboard();
  } catch (error: unknown) {
    dashboard.value = null;
    dashboardErrorMessage.value = getDashboardErrorMessage(error);
  } finally {
    isDashboardLoading.value = false;
  }
};

const getLinksErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 400) {
      return "Проверьте параметры поиска списка ссылок.";
    }

    if (error.status === 401) {
      return "Сессия недействительна. Войдите заново, чтобы открыть список ссылок.";
    }
  }

  return "Не удалось загрузить список ссылок. Проверьте соединение и повторите попытку.";
};

const loadLinks = async () => {
  const requestId = ++linksRequestId;

  isLinksLoading.value = true;
  linksErrorMessage.value = "";

  try {
    const response = await listLinks({
      page: linksPage.value,
      pageSize: linksPageSize,
      q: linksQ.value,
    });

    if (requestId !== linksRequestId) {
      return;
    }

    links.value = response.items;
    linksPagination.value = response.pagination;
  } catch (error: unknown) {
    if (requestId !== linksRequestId) {
      return;
    }

    links.value = [];
    linksPagination.value = null;
    linksErrorMessage.value = getLinksErrorMessage(error);
  } finally {
    if (requestId === linksRequestId) {
      isLinksLoading.value = false;
    }
  }
};

const onLinkFiltersChange = (filters: { q: string }) => {
  linksQ.value = filters.q;
  linksPage.value = 1;
  void loadLinks();
};

const onLinkUpdated = (updatedLink: Link) => {
  links.value = links.value.map((link) => (link.id === updatedLink.id ? updatedLink : link));
  void loadDashboard();
};

const onLinksPageChange = (page: number) => {
  linksPage.value = page;
  void loadLinks();
};

const goToCreateLink = () => {
  void router.push(ROUTES.linkCreate);
};

onMounted(() => {
  void loadDashboard();
  void loadLinks();
});
</script>

<template>
  <section class="dashboard-page">
    <UiPageHeader title="Dashboard" subtitle="Ключевые показатели аккаунта и последние созданные ссылки." />

    <DashboardSummary
      v-if="isDashboardLoading || dashboard"
      :total-links="summary.totalLinks"
      :active-links="summary.activeLinks"
      :total-clicks="summary.totalClicks"
      :clicks-last24h="summary.clicksLast24h"
      :loading="isDashboardLoading"
    />

    <UiPageState
      v-if="isDashboardLoading && !dashboard"
      type="loading"
      title="Загружаем dashboard"
      description="Получаем метрики и последние ссылки аккаунта."
    />

    <UiPageState
      v-else-if="dashboardErrorMessage"
      type="error"
      title="Dashboard недоступен"
      :description="dashboardErrorMessage"
      action-text="Повторить"
      @action="loadDashboard"
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

      <LinksTable
        :links="links"
        :pagination="linksPagination"
        :loading="isLinksLoading"
        :error-message="linksErrorMessage"
        :q="linksQ"
        @filters-change="onLinkFiltersChange"
        @link-updated="onLinkUpdated"
        @page-change="onLinksPageChange"
        @retry="loadLinks"
      />
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
