<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { listAdminLinks, type AdminLink } from "@/api/admin";
import type { ApiClientError } from "@/api/types";
import type { Pagination } from "@/entities/link/link.types";
import { useSession } from "@/entities/session/useSession";
import { AdminLinksSearch, AdminLinksTable } from "@/features/admin-link-moderation";
import { useI18n } from "@/shared/composables/useI18n";
import { UiPageHeader, UiPageState } from "@/shared/ui";

const session = useSession();
const { t } = useI18n();

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
    return t("admin.error.access");
  }

  if (isApiClientError(error) && error.status === 400) {
    return t("admin.error.badSearch");
  }

  return t("admin.error.failed");
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
      :title="t('admin.page.checkingTitle')"
      :description="t('admin.page.checkingDescription')"
    />

    <UiPageState
      v-else-if="isForbidden"
      type="forbidden"
      :title="t('admin.page.forbiddenTitle')"
      :description="t('admin.page.forbiddenDescription')"
    />

    <UiPageState
      v-else-if="accessErrorMessage"
      type="forbidden"
      :title="t('admin.page.forbiddenTitle')"
      :description="accessErrorMessage"
    />

    <template v-else>
      <UiPageHeader
        :title="t('admin.page.title')"
        :subtitle="t('admin.page.subtitle')"
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
