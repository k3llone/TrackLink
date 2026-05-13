<script setup lang="ts">
import { UiButton, UiSwitch } from "@/shared/ui";
</script>

<template>
  <section class="security-actions" aria-labelledby="security-actions-title">
    <header class="security-actions__header">
      <span class="security-actions__icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" focusable="false">
          <path
            d="M12 3 5 6v5c0 4.4 2.8 8.4 7 10 4.2-1.6 7-5.6 7-10V6l-7-3Zm0 2.2 5 2.1V11c0 3.3-2 6.4-5 7.7-3-1.3-5-4.4-5-7.7V7.3l5-2.1Z"
            fill="currentColor"
          />
        </svg>
      </span>

      <div class="security-actions__title-group">
        <h2 id="security-actions-title" class="security-actions__title">Security Settings</h2>
        <p class="security-actions__subtitle">Manage your account security</p>
      </div>
    </header>

    <div class="security-actions__list">
      <div class="security-actions__row">
        <div class="security-actions__copy">
          <div class="security-actions__row-heading">
            <h3 class="security-actions__row-title">Two-Factor Authentication</h3>
            <span class="security-actions__tag">Post-MVP</span>
          </div>
          <p class="security-actions__row-subtitle">Add an extra layer of security to your account</p>
        </div>

        <UiSwitch
          class="security-actions__switch"
          :model-value="false"
          disabled
          aria-label="Two-factor authentication is planned for a later release"
        />
      </div>

      <div class="security-actions__sessions">
        <div class="security-actions__copy">
          <h3 class="security-actions__row-title">Active Sessions</h3>
          <p class="security-actions__row-subtitle">Manage your active login sessions</p>
        </div>

        <div class="security-actions__row security-actions__row--session">
          <div class="security-actions__copy">
            <div class="security-actions__row-heading">
              <h4 class="security-actions__session-title">Current Session</h4>
              <span class="security-actions__tag">Post-MVP</span>
            </div>
            <p class="security-actions__row-subtitle">Last active: Just now</p>
          </div>

          <UiButton
            class="security-actions__deactivate"
            size="sm"
            disabled
            title="Session management is planned for a later release"
          >
            Deactivate
          </UiButton>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.security-actions {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.security-actions__header {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.security-actions__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  margin-top: 2px;
  color: var(--tl-color-primary);
}

.security-actions__icon svg {
  width: 20px;
  height: 20px;
}

.security-actions__title-group,
.security-actions__copy {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.security-actions__title {
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.security-actions__subtitle,
.security-actions__row-subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.security-actions__list,
.security-actions__sessions {
  display: flex;
  flex-direction: column;
  gap: 18px;
  width: min(100%, 950px);
}

.security-actions__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  min-height: 54px;
}

.security-actions__row--session {
  padding-top: 8px;
}

.security-actions__row-heading {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.security-actions__row-title,
.security-actions__session-title {
  color: var(--tl-color-text);
  font-size: 16px;
  line-height: 1.3;
}

.security-actions__tag {
  display: inline-flex;
  align-items: center;
  min-height: 22px;
  border-radius: 999px;
  padding: 3px 8px;
  background: var(--tl-color-surface-muted);
  color: var(--tl-color-text-muted);
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
}

.security-actions__switch {
  width: 44px;
  min-width: 44px;
  padding: 0;
}

.security-actions__switch :deep(.ui-switch__content) {
  display: none;
}

.security-actions__deactivate {
  width: 180px;
}

@media (max-width: 767px) {
  .security-actions__row {
    align-items: flex-start;
    flex-direction: column;
  }

  .security-actions__deactivate {
    width: 100%;
  }
}
</style>
