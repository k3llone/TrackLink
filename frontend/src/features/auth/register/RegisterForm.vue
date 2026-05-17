<script setup lang="ts">
import { RouterLink } from "vue-router";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton, UiInput } from "@/shared/ui";
import { ROUTES } from "@/shared/lib/routes/paths";
import AuthFormCard from "../AuthFormCard.vue";
import { useRegisterForm } from "./useRegisterForm";

const { form, errors, isSubmitting, submit } = useRegisterForm();
const { t } = useI18n();
</script>

<template>
  <AuthFormCard :title="t('auth.register.title')" :subtitle="t('auth.register.subtitle')">
    <form class="register-form" novalidate @submit.prevent="submit">
      <p v-if="errors.form" class="register-form__error" role="alert">{{ errors.form }}</p>

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
        autocomplete="new-password"
        :placeholder="t('auth.register.passwordPlaceholder')"
        required
        :disabled="isSubmitting"
        :error="errors.password"
      />

      <UiInput
        v-model="form.confirmPassword"
        :label="t('auth.register.confirmPasswordLabel')"
        type="password"
        autocomplete="new-password"
        :placeholder="t('auth.register.confirmPasswordPlaceholder')"
        required
        :disabled="isSubmitting"
        :error="errors.confirmPassword"
      />

      <UiButton type="submit" size="lg" full-width :loading="isSubmitting">{{ t("auth.register.submit") }}</UiButton>
    </form>

    <template #footer>
      {{ t("auth.register.footerText") }}
      <RouterLink class="register-form__footer-link" :to="ROUTES.login">{{ t("auth.register.footerAction") }}</RouterLink>
    </template>
  </AuthFormCard>
</template>

<style scoped>
.register-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.register-form__error {
  border-radius: var(--tl-radius-md);
  padding: 10px 12px;
  background: rgb(225 75 90 / 10%);
  color: var(--tl-color-danger);
  font-size: 13px;
}

.register-form__footer-link {
  color: var(--tl-color-primary);
  font-weight: 700;
  text-decoration: none;
}

.register-form__footer-link:hover {
  text-decoration: underline;
}
</style>
