<script setup lang="ts">
import { computed, onMounted } from "vue";
import ProfileInfoCard from "@/features/account-settings/ProfileInfoCard.vue";
import SecurityActions from "@/features/account-settings/SecurityActions.vue";
import PostMvpSettingsSection from "@/features/account-settings/PostMvpSettingsSection.vue";
import { useSession } from "@/entities/session/useSession";
import { UiPageHeader, UiPageState } from "@/shared/ui";

const session = useSession();

const currentUser = computed(() => session.user.value);
const isLoading = computed(() => session.isLoading.value);

onMounted(() => {
  if (session.status.value === "idle") {
    void session.loadCurrentUser();
  }
});
</script>

<template>
  <section class="settings-page">
    <UiPageHeader title="Account settings" subtitle="Manage profile details and account security." />

    <UiPageState
      v-if="isLoading"
      type="loading"
      title="Loading account"
      description="Preparing your account settings."
    />

    <UiPageState
      v-else-if="!currentUser"
      type="error"
      title="Account data unavailable"
      description="Refresh the page or sign in again to load your account details."
    />

    <div v-else class="settings-page__content">
      <ProfileInfoCard :email="currentUser.email" :role="currentUser.role" :created-at="currentUser.createdAt" />
      <PostMvpSettingsSection />
      <SecurityActions />
    </div>
  </section>
</template>

<style scoped>
.settings-page {
  --settings-section-divider: rgb(37 31 63 / 10%);

  width: 100%;
}

.settings-page__content {
  display: flex;
  flex-direction: column;
  gap: 28px;
  width: 100%;
  max-width: 1020px;
}

.settings-page__content > * + * {
  border-top: 1px solid var(--settings-section-divider);
  padding-top: 28px;
}

@media (max-width: 767px) {
  .settings-page__content {
    gap: 24px;
  }

  .settings-page__content > * + * {
    padding-top: 24px;
  }
}
</style>
