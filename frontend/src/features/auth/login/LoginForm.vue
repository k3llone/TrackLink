<script setup lang="ts">
import { RouterLink } from "vue-router";
import { UiButton, UiInput } from "@/shared/ui";
import { ROUTES } from "@/shared/lib/routes/paths";
import { useLoginForm } from "./useLoginForm";

const { form, errors, isSubmitting, submit } = useLoginForm();
</script>

<template>
  <section class="login-form" aria-labelledby="login-title">
    <div class="login-form__header">
      <h1 id="login-title" class="login-form__title">Welcome back</h1>
      <p class="login-form__subtitle">Sign in to manage your TrackLink links.</p>
    </div>

    <form class="login-form__body" novalidate @submit.prevent="submit">
      <p v-if="errors.form" class="login-form__error" role="alert">{{ errors.form }}</p>

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
        autocomplete="current-password"
        placeholder="Enter your password"
        required
        :disabled="isSubmitting"
        :error="errors.password"
      />

      <RouterLink class="login-form__forgot" :to="ROUTES.forgotPassword">Forgot password?</RouterLink>

      <UiButton type="submit" size="lg" full-width :loading="isSubmitting">Sign in</UiButton>
    </form>

    <p class="login-form__footer">
      Don't have an account?
      <RouterLink :to="ROUTES.register">Create account</RouterLink>
    </p>
  </section>
</template>

<style scoped>
.login-form {
  width: 100%;
  padding: 32px;
  border: 1px solid #e5deef;
  border-radius: 24px;
  background: var(--tl-color-white);
  box-shadow: 0 18px 48px rgb(37 31 63 / 10%);
}

.login-form__header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 28px;
  text-align: center;
}

.login-form__title {
  color: var(--tl-color-text);
  font-size: 28px;
  line-height: 1.2;
}

.login-form__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.login-form__body {
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

.login-form__forgot:hover {
  text-decoration: underline;
}

.login-form__footer {
  margin-top: 22px;
  color: var(--tl-color-text-muted);
  font-size: 14px;
  text-align: center;
}

.login-form__footer a {
  color: var(--tl-color-primary);
  font-weight: 700;
  text-decoration: none;
}

.login-form__footer a:hover {
  text-decoration: underline;
}

@media (max-width: 767px) {
  .login-form {
    padding: 24px;
  }
}
</style>
