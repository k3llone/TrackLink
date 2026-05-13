<script setup lang="ts">
import { UiButton, UiStatusBadge } from "@/shared/ui";
</script>

<template>
  <section class="post-mvp-settings" aria-labelledby="post-mvp-settings-title">
    <div class="post-mvp-settings__copy">
      <div class="post-mvp-settings__heading">
        <h2 id="post-mvp-settings-title" class="post-mvp-settings__title">Account Tier</h2>
        <UiStatusBadge status="pending" label="Post-MVP" />
      </div>
      <p class="post-mvp-settings__subtitle">Plan and limit management will be connected after the backend contract.</p>
    </div>

    <div class="post-mvp-settings__actions">
      <span class="post-mvp-settings__tier">Free</span>
      <UiButton size="sm" disabled title="Upgrade flow is planned for a later release">
        Upgrade
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
