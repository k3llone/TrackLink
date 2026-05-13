<script setup lang="ts">
import { computed } from "vue";
import type { UserRole } from "@/api/auth";
import { UiButton, UiInput, UiStatusBadge } from "@/shared/ui";

const props = defineProps<{
  email: string;
  role: UserRole;
  createdAt?: string;
}>();

const maskedPassword = "********";

const formattedCreatedAt = computed(() => {
  if (!props.createdAt) {
    return "";
  }

  const createdDate = new Date(props.createdAt);

  if (Number.isNaN(createdDate.getTime())) {
    return "";
  }

  return new Intl.DateTimeFormat("en", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(createdDate);
});

const roleBadgeStatus = computed<"active" | "inactive">(() => (props.role === "admin" ? "active" : "inactive"));
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
        <h2 id="profile-info-title" class="profile-info__title">Profile Information</h2>
        <p class="profile-info__subtitle">Manage your account details</p>
      </div>
    </header>

    <div class="profile-info__meta" aria-label="Account metadata">
      <div class="profile-info__meta-item">
        <span class="profile-info__meta-label">Role</span>
        <UiStatusBadge :status="roleBadgeStatus" :label="role" />
      </div>

      <div v-if="formattedCreatedAt" class="profile-info__meta-item">
        <span class="profile-info__meta-label">Created</span>
        <strong class="profile-info__meta-value">{{ formattedCreatedAt }}</strong>
      </div>
    </div>

    <div class="profile-info__fields">
      <UiInput :model-value="email" label="Email" type="email" autocomplete="email" readonly />

      <UiInput
        :model-value="maskedPassword"
        label="Password"
        type="text"
        autocomplete="current-password"
        readonly
      />

      <div class="profile-info__action">
        <UiButton size="sm" disabled title="Password change is planned for a later auth feature">
          Change password
        </UiButton>
        <p class="profile-info__hint">Password change will be available after the dedicated auth flow is implemented.</p>
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
