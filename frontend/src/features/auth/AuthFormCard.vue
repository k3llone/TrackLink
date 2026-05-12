<script setup lang="ts">
defineProps<{
  title: string;
  subtitle: string;
}>();
</script>

<template>
  <section class="auth-form-card" aria-labelledby="auth-form-title">
    <div class="auth-form-card__header">
      <h1 id="auth-form-title" class="auth-form-card__title">{{ title }}</h1>
      <p class="auth-form-card__subtitle">{{ subtitle }}</p>
    </div>

    <div class="auth-form-card__body">
      <slot />
    </div>

    <div v-if="$slots.footer" class="auth-form-card__footer">
      <slot name="footer" />
    </div>
  </section>
</template>

<style scoped>
.auth-form-card {
  width: 100%;
  padding: 32px;
  border: 1px solid #e5deef;
  border-radius: 24px;
  background: var(--tl-color-white);
  box-shadow: 0 18px 48px rgb(37 31 63 / 10%);
}

.auth-form-card__header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 28px;
  text-align: center;
}

.auth-form-card__title {
  color: var(--tl-color-text);
  font-size: 28px;
  line-height: 1.2;
}

.auth-form-card__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.auth-form-card__body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.auth-form-card__footer {
  margin-top: 22px;
  color: var(--tl-color-text-muted);
  font-size: 14px;
  text-align: center;
}

@media (max-width: 767px) {
  .auth-form-card {
    padding: 24px;
  }
}
</style>
