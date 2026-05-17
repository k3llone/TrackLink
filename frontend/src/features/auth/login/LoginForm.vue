<script setup lang="ts">
import { RouterLink } from "vue-router";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton, UiInput } from "@/shared/ui";
import { ROUTES } from "@/shared/lib/routes/paths";
import AuthFormCard from "../AuthFormCard.vue";
import { useLoginForm } from "./useLoginForm";

const { form, errors, isSubmitting, submit } = useLoginForm();
const { t } = useI18n();
</script>

<template>
  <AuthFormCard :title="t('auth.login.title')" :subtitle="t('auth.login.subtitle')">
    <form class="login-form" novalidate @submit.prevent="submit">
      <p v-if="errors.form" class="login-form__error" role="alert">{{ errors.form }}</p>

      <UiInput
        v-model="form.email"
        :label="t('common.email')"
        type="email"
        autocomplete="email"
        placeholder="you@example.com"
        required
        :disabled="isSubmitting"
        :error="errors.email"
      />

      <UiInput
        v-model="form.password"
        :label="t('common.password')"
        type="password"
        autocomplete="current-password"
        :placeholder="t('auth.login.passwordPlaceholder')"
        required
        :disabled="isSubmitting"
        :error="errors.password"
      />

      <RouterLink class="login-form__forgot" :to="ROUTES.forgotPassword">{{ t("auth.login.forgotPassword") }}</RouterLink>

      <UiButton type="submit" size="lg" full-width :loading="isSubmitting">{{ t("auth.login.submit") }}</UiButton>
    </form>

    <template #footer>
      {{ t("auth.login.footerText") }}
      <RouterLink class="login-form__footer-link" :to="ROUTES.register">{{ t("auth.login.footerAction") }}</RouterLink>
    </template>
  </AuthFormCard>
</template>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.login-form__error {
  border-radius: var(--tl-radius-md);
  padding: 10px 12px;
  background: rgb(225 75 90 / 10%);
  color: var(--tl-color-danger);
  font-size: 13px;
}

.login-form__forgot {
  align-self: flex-end;
  color: var(--tl-color-primary);
  font-size: 13px;
  font-weight: 600;
  text-decoration: none;
}

.login-form__forgot:hover,
.login-form__footer-link:hover {
  text-decoration: underline;
}

.login-form__footer-link {
  color: var(--tl-color-primary);
  font-weight: 700;
  text-decoration: none;
}
</style>
