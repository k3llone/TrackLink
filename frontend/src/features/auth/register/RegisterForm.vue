<script setup lang="ts">
import { RouterLink } from "vue-router";
import { UiButton, UiInput } from "@/shared/ui";
import { ROUTES } from "@/shared/lib/routes/paths";
import AuthFormCard from "../AuthFormCard.vue";
import { useRegisterForm } from "./useRegisterForm";

const { form, errors, isSubmitting, submit } = useRegisterForm();
</script>

<template>
  <AuthFormCard title="Create your account" subtitle="Start managing short links with TrackLink.">
    <form class="register-form" novalidate @submit.prevent="submit">
      <p v-if="errors.form" class="register-form__error" role="alert">{{ errors.form }}</p>

      <UiInput
        v-model="form.email"
        label="Email"
        type="email"
        autocomplete="email"
        placeholder="you@example.com"
        required
        :disabled="isSubmitting"
        :error="errors.email"
      />

      <UiInput
        v-model="form.password"
        label="Password"
        type="password"
        autocomplete="new-password"
        placeholder="At least 8 characters"
        required
        :disabled="isSubmitting"
        :error="errors.password"
      />

      <UiInput
        v-model="form.confirmPassword"
        label="Confirm password"
        type="password"
        autocomplete="new-password"
        placeholder="Repeat your password"
        required
        :disabled="isSubmitting"
        :error="errors.confirmPassword"
      />

      <UiButton type="submit" size="lg" full-width :loading="isSubmitting">Create account</UiButton>
    </form>

    <template #footer>
      Already have an account?
      <RouterLink class="register-form__footer-link" :to="ROUTES.login">Sign in</RouterLink>
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
