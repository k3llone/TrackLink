<script setup lang="ts">
import { UiBrandLogo } from "@/shared/ui";
</script>

<template>
  <main class="auth-layout">
    <section class="auth-layout__content">
      <UiBrandLogo variant="auth" />
      <div class="auth-layout__slot">
        <slot />
      </div>
    </section>
  </main>
</template>

<style scoped>
.auth-layout {
  min-height: 100dvh;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 40px 16px;
  background: var(--tl-color-bg);
}

.auth-layout__content {
  width: var(--tl-form-width);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 28px;
}

.auth-layout__slot {
  width: 100%;
}

@media (max-width: 1023px) {
  .auth-layout {
    padding-inline: 24px;
  }

  .auth-layout__content {
    width: min(var(--tl-form-width), calc(100vw - 48px));
  }
}

@media (max-width: 767px) {
  .auth-layout {
    padding-inline: 16px;
  }

  .auth-layout__content {
    width: calc(100vw - 32px);
  }
}
</style>
