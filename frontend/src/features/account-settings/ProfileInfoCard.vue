<script setup lang="ts">
import { computed } from "vue";
import type { UserRole } from "@/api/auth";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton, UiInput, UiStatusBadge } from "@/shared/ui";

const props = defineProps<{
  email: string;
  role: UserRole;
  createdAt?: string;
}>();

const { formatDate, t } = useI18n();
const maskedPassword = "********";

const formattedCreatedAt = computed(() => {
  if (!props.createdAt) {
    return "";
  }

  return formatDate(props.createdAt);
});

const roleBadgeStatus = computed<"active" | "inactive">(() => (props.role === "admin" ? "active" : "inactive"));
const roleLabel = computed(() => (props.role === "admin" ? t("user.role.admin") : t("user.role.customer")));
</script>

<template>
  <section class="profile-info" aria-labelledby="profile-info-title">
    <header class="profile-info__header">
      <span class="profile-info__icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" focusable="false">
          <path
            d="M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Zm0 2c-4 0-7 2-7 4.5V20h14v-1.5C19 16 16 14 12 14Z"
            fill="currentColor"
          />
        </svg>
      </span>

      <div class="profile-info__title-group">
        <h2 id="profile-info-title" class="profile-info__title">{{ t("settings.profile.title") }}</h2>
        <p class="profile-info__subtitle">{{ t("settings.profile.subtitle") }}</p>
      </div>
    </header>

    <div class="profile-info__meta" :aria-label="t('settings.profile.metadataAria')">
      <div class="profile-info__meta-item">
        <span class="profile-info__meta-label">{{ t("settings.profile.role") }}</span>
        <UiStatusBadge :status="roleBadgeStatus" :label="roleLabel" />
      </div>

      <div v-if="formattedCreatedAt" class="profile-info__meta-item">
        <span class="profile-info__meta-label">{{ t("settings.profile.created") }}</span>
        <strong class="profile-info__meta-value">{{ formattedCreatedAt }}</strong>
      </div>
    </div>

    <div class="profile-info__fields">
      <UiInput :model-value="email" :label="t('common.email')" type="email" autocomplete="email" readonly />

      <UiInput
        :model-value="maskedPassword"
        :label="t('common.password')"
        type="text"
        autocomplete="current-password"
        readonly
      />

      <div class="profile-info__action">
        <UiButton size="sm" disabled :title="t('settings.profile.changePasswordTitle')">
          {{ t("settings.profile.changePassword") }}
        </UiButton>
        <p class="profile-info__hint">{{ t("settings.profile.changePasswordHint") }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.profile-info {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.profile-info__header {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.profile-info__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  margin-top: 2px;
  color: var(--tl-color-primary);
}

.profile-info__icon svg {
  width: 20px;
  height: 20px;
}

.profile-info__title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.profile-info__title {
  color: var(--tl-color-text);
  font-size: 20px;
  line-height: 1.25;
}

.profile-info__subtitle,
.profile-info__hint {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.profile-info__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.profile-info__meta-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 32px;
}

.profile-info__meta-label {
  color: var(--tl-color-text-muted);
  font-size: 13px;
  font-weight: 700;
}

.profile-info__meta-value {
  color: var(--tl-color-text);
  font-size: 14px;
}

.profile-info__fields {
  display: grid;
  gap: 16px;
  width: min(100%, 950px);
}

.profile-info__action {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.profile-info__hint {
  max-width: 520px;
}

@media (max-width: 767px) {
  .profile-info__action {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
