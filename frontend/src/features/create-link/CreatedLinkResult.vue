<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { Link, LinkStatus } from "@/entities/link/link.types";
import CopyShortUrlButton from "@/features/link-actions/CopyShortUrlButton.vue";
import { getLinkDetailsPath, ROUTES } from "@/shared/lib/routes/paths";
import { UiButton, UiStatusBadge } from "@/shared/ui";

const props = defineProps<{
  link: Link;
}>();

const emit = defineEmits<{
  "create-more": [];
}>();

const router = useRouter();
const numberFormatter = new Intl.NumberFormat("ru-RU");

const statusLabels: Record<LinkStatus, string> = {
  active: "Активна",
  inactive: "Неактивна",
  blocked: "Заблокирована",
  deleted: "Удалена",
};

const shortUrl = computed(() => props.link.shortUrl || props.link.code);
const totalClicksLabel = computed(() => numberFormatter.format(props.link.totalClicks));
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
        <h2 id="created-link-result-title" class="created-link-result__title">Короткая ссылка создана</h2>
        <p class="created-link-result__subtitle">Short URL готов к использованию и уже сохранён в вашем аккаунте.</p>
      </div>
    </header>

    <div class="created-link-result__short-url">
      <span class="created-link-result__label">Short URL</span>
      <a class="created-link-result__url" :href="shortUrl" target="_blank" rel="noreferrer" :title="shortUrl">
        {{ shortUrl }}
      </a>
      <CopyShortUrlButton :short-url="shortUrl" variant="secondary" size="md" />
    </div>

    <dl class="created-link-result__details">
      <div class="created-link-result__detail">
        <dt>Target URL</dt>
        <dd>
          <a :href="link.targetUrl" target="_blank" rel="noreferrer" :title="link.targetUrl">
            {{ link.targetUrl }}
          </a>
        </dd>
      </div>

      <div class="created-link-result__detail">
        <dt>Статус</dt>
        <dd>
          <UiStatusBadge :status="link.status" :label="statusLabels[link.status]" />
        </dd>
      </div>

      <div v-if="shouldShowTotalClicks" class="created-link-result__detail">
        <dt>Переходы</dt>
        <dd>{{ totalClicksLabel }}</dd>
      </div>
    </dl>

    <footer class="created-link-result__actions">
      <UiButton type="button" @click="openAnalytics">Открыть аналитику</UiButton>
      <UiButton variant="secondary" type="button" @click="createMore">Создать ещё</UiButton>
      <UiButton variant="ghost" type="button" @click="goToDashboard">На dashboard</UiButton>
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
