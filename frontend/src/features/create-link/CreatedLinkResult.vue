<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { Link, LinkStatus } from "@/entities/link/link.types";
import CopyShortUrlButton from "@/features/link-actions/CopyShortUrlButton.vue";
import { useI18n } from "@/shared/composables/useI18n";
import { getLinkDetailsPath, ROUTES } from "@/shared/lib/routes/paths";
import { UiButton, UiStatusBadge } from "@/shared/ui";

const props = defineProps<{
  link: Link;
}>();

const emit = defineEmits<{
  "create-more": [];
}>();

const router = useRouter();
const { formatNumber, t } = useI18n();

const statusLabels: Record<LinkStatus, () => string> = {
  active: () => t("link.status.active"),
  inactive: () => t("link.status.inactive"),
  blocked: () => t("link.status.blocked"),
  deleted: () => t("link.status.deleted"),
};

const shortUrl = computed(() => props.link.shortUrl || props.link.code);
const totalClicksLabel = computed(() => formatNumber(props.link.totalClicks));
const shouldShowTotalClicks = computed(() => typeof props.link.totalClicks === "number");

const openAnalytics = () => {
  void router.push(getLinkDetailsPath(props.link.id));
};

const createMore = () => {
  emit("create-more");
};

const goToDashboard = () => {
  void router.push(ROUTES.dashboard);
};
</script>

<template>
  <section class="created-link-result" aria-labelledby="created-link-result-title">
    <header class="created-link-result__header">
      <div class="created-link-result__title-group">
        <h2 id="created-link-result-title" class="created-link-result__title">{{ t("createLink.result.title") }}</h2>
        <p class="created-link-result__subtitle">{{ t("createLink.result.subtitle") }}</p>
      </div>
    </header>

    <div class="created-link-result__short-url">
      <span class="created-link-result__label">{{ t("common.shortUrl") }}</span>
      <a class="created-link-result__url" :href="shortUrl" target="_blank" rel="noreferrer" :title="shortUrl">
        {{ shortUrl }}
      </a>
      <CopyShortUrlButton :short-url="shortUrl" variant="secondary" size="md" />
    </div>

    <dl class="created-link-result__details">
      <div class="created-link-result__detail">
        <dt>{{ t("common.targetUrl") }}</dt>
        <dd>
          <a :href="link.targetUrl" target="_blank" rel="noreferrer" :title="link.targetUrl">
            {{ link.targetUrl }}
          </a>
        </dd>
      </div>

      <div class="created-link-result__detail">
        <dt>{{ t("common.status") }}</dt>
        <dd>
          <UiStatusBadge :status="link.status" :label="statusLabels[link.status]()" />
        </dd>
      </div>

      <div v-if="shouldShowTotalClicks" class="created-link-result__detail">
        <dt>{{ t("common.clicks") }}</dt>
        <dd>{{ totalClicksLabel }}</dd>
      </div>
    </dl>

    <footer class="created-link-result__actions">
      <UiButton type="button" @click="openAnalytics">{{ t("createLink.result.openAnalytics") }}</UiButton>
      <UiButton variant="secondary" type="button" @click="createMore">{{ t("createLink.result.createMore") }}</UiButton>
      <UiButton variant="ghost" type="button" @click="goToDashboard">{{ t("createLink.result.toDashboard") }}</UiButton>
    </footer>
  </section>
</template>

<style scoped>
.created-link-result {
  display: flex;
  flex-direction: column;
  gap: 18px;
  width: min(100%, 720px);
  padding: 22px;
  border: 1px solid rgb(37 31 63 / 10%);
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-white);
}

.created-link-result__header,
.created-link-result__title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.created-link-result__title {
  color: var(--tl-color-text);
  font-size: 22px;
  line-height: 1.25;
}

.created-link-result__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.created-link-result__short-url {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px 12px;
  align-items: center;
  padding: 16px;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
}

.created-link-result__label {
  grid-column: 1 / -1;
  color: var(--tl-color-text-muted);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.created-link-result__url {
  min-width: 0;
  color: var(--tl-color-primary);
  font-size: 18px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  text-decoration: none;
  white-space: nowrap;
}

.created-link-result__url:hover,
.created-link-result__detail a:hover {
  text-decoration: underline;
}

.created-link-result__details {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

.created-link-result__detail {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
  padding: 14px;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
}

.created-link-result__detail dt {
  color: var(--tl-color-text-muted);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.created-link-result__detail dd {
  min-width: 0;
  margin: 0;
  color: var(--tl-color-text);
  font-size: 14px;
  font-weight: 600;
}

.created-link-result__detail a {
  display: inline-block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  text-decoration: none;
  vertical-align: bottom;
  white-space: nowrap;
}

.created-link-result__actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

@media (max-width: 767px) {
  .created-link-result {
    padding: 18px;
  }

  .created-link-result__short-url,
  .created-link-result__details {
    grid-template-columns: 1fr;
  }

  .created-link-result__actions,
  .created-link-result__actions :deep(.ui-button) {
    width: 100%;
  }
}
</style>
