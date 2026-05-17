<script setup lang="ts">
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton, UiStatusBadge } from "@/shared/ui";

const { t } = useI18n();
</script>

<template>
  <section class="post-mvp-settings" aria-labelledby="post-mvp-settings-title">
    <div class="post-mvp-settings__copy">
      <div class="post-mvp-settings__heading">
        <h2 id="post-mvp-settings-title" class="post-mvp-settings__title">{{ t("settings.tier.title") }}</h2>
        <UiStatusBadge status="pending" :label="t('common.postMvp')" />
      </div>
      <p class="post-mvp-settings__subtitle">{{ t("settings.tier.subtitle") }}</p>
    </div>

    <div class="post-mvp-settings__actions">
      <span class="post-mvp-settings__tier">{{ t("common.free") }}</span>
      <UiButton size="sm" disabled :title="t('settings.tier.upgradeTitle')">
        {{ t("common.upgrade") }}
      </UiButton>
    </div>
  </section>
</template>

<style scoped>
.post-mvp-settings {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  width: min(100%, 950px);
}

.post-mvp-settings__copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.post-mvp-settings__heading {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.post-mvp-settings__title {
  color: var(--tl-color-text);
  font-size: 16px;
  line-height: 1.3;
}

.post-mvp-settings__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.post-mvp-settings__actions {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.post-mvp-settings__tier {
  color: var(--tl-color-primary);
  font-size: 16px;
  font-weight: 800;
}

@media (max-width: 767px) {
  .post-mvp-settings {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
