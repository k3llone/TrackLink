<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { listAdminLinks, type AdminLink } from "@/api/admin";
import type { ApiClientError } from "@/api/types";
import type { Pagination } from "@/entities/link/link.types";
import { useSession } from "@/entities/session/useSession";
import { AdminLinksSearch, AdminLinksTable } from "@/features/admin-link-moderation";
import { UiPageHeader, UiPageState } from "@/shared/ui";

const session = useSession();

const links = ref<AdminLink[]>([]);
const linksPagination = ref<Pagination | null>(null);
const isLinksLoading = ref(false);
const linksErrorMessage = ref("");
const accessErrorMessage = ref("");
const linksPage = ref(1);
const linksPageSize = 20;
const linksQ = ref("");
let linksRequestId = 0;

const isCheckingAccess = computed(() => session.isSessionLoading.value && !session.isAdmin.value);
const isForbidden = computed(() => !isCheckingAccess.value && !session.isAdmin.value);
const hasLinkSearch = computed(() => Boolean(linksQ.value.trim()));

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const isAccessError = (error: unknown) =>
  isApiClientError(error) && (error.status === 401 || error.status === 403);

const getLinksErrorMessage = (error: unknown) => {
  if (isAccessError(error)) {
    return "У вас нет доступа к административной панели.";
  }

  if (isApiClientError(error) && error.status === 400) {
    return "Проверьте параметры поиска ссылок.";
  }

  return "Не удалось загрузить список ссылок.";
};

const loadLinks = async () => {
  if (!session.isAdmin.value) {
    return;
  }

  const requestId = ++linksRequestId;

  isLinksLoading.value = true;
  linksErrorMessage.value = "";
  accessErrorMessage.value = "";

  try {
    const response = await listAdminLinks({
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

    const message = getLinksErrorMessage(error);
    if (isAccessError(error)) {
      accessErrorMessage.value = message;
    } else {
      linksErrorMessage.value = message;
    }
  } finally {
    if (requestId === linksRequestId) {
      isLinksLoading.value = false;
    }
  }
};

const initialize = async () => {
  await session.loadCurrentUser();

  if (session.isAdmin.value) {
    await loadLinks();
  }
};

const onLinksPageChange = (page: number) => {
  linksPage.value = page;
  void loadLinks();
};

const onLinksSearchChange = (q: string) => {
  linksQ.value = q;
  linksPage.value = 1;
  void loadLinks();
};

const onLinkUpdated = (updatedLink: AdminLink) => {
  links.value = links.value.map((link) =>
    link.id === updatedLink.id ? { ...link, ...updatedLink } : link,
  );
};

onMounted(() => {
  void initialize();
});
</script>

<template>
  <section class="admin-page">
    <UiPageState
      v-if="isCheckingAccess"
      type="loading"
      title="Проверяем доступ"
      description="Проверяем права доступа к административной панели."
    />

    <UiPageState
      v-else-if="isForbidden"
      type="forbidden"
      title="Нет доступа"
      description="У вас нет доступа к административной панели."
    />

    <UiPageState
      v-else-if="accessErrorMessage"
      type="forbidden"
      title="Нет доступа"
      :description="accessErrorMessage"
    />

    <template v-else>
      <UiPageHeader
        title="Admin panel"
        subtitle="Просмотр коротких ссылок и административная блокировка."
      />

      <div class="admin-page__content">
        <AdminLinksSearch :q="linksQ" @change="onLinksSearchChange" />

        <AdminLinksTable
          :links="links"
          :pagination="linksPagination"
          :loading="isLinksLoading"
          :error-message="linksErrorMessage"
          :has-search="hasLinkSearch"
          @page-change="onLinksPageChange"
          @link-updated="onLinkUpdated"
          @retry="loadLinks"
        />
      </div>
    </template>
  </section>
</template>

<style scoped>
.admin-page {
  width: 100%;
}

.admin-page__content {
  display: flex;
  flex-direction: column;
  gap: 18px;
}
</style>
