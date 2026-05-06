<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: "auth" | "header";
    showText?: boolean;
  }>(),
  {
    variant: "header",
    showText: true,
  },
);
</script>

<template>
  <div class="ui-brand-logo" :class="`ui-brand-logo--${variant}`">
    <span class="ui-brand-logo__icon" aria-hidden="true">
      <svg viewBox="0 0 24 24" focusable="false">
        <path
          d="M8.3 7.2a3 3 0 0 1 4.2 0l1.5 1.5a3 3 0 0 1 0 4.2l-1.1 1.1-1.4-1.4 1.1-1.1a1 1 0 0 0 0-1.4L11 8.6a1 1 0 0 0-1.4 0L8.5 9.8 7.1 8.4l1.2-1.2Zm7.6 9.6a3 3 0 0 1-4.2 0l-1.5-1.5a3 3 0 0 1 0-4.2l1.1-1.1 1.4 1.4-1.1 1.1a1 1 0 0 0 0 1.4l1.5 1.5a1 1 0 0 0 1.4 0l1.1-1.1 1.4 1.4-1.1 1.1ZM8 12h8v2H8v-2Z"
          fill="currentColor"
        />
      </svg>
    </span>
    <span v-if="showText" class="ui-brand-logo__text">TrackLink</span>
  </div>
</template>

<style scoped>
.ui-brand-logo {
  display: inline-flex;
  align-items: center;
  color: var(--tl-color-primary);
}

.ui-brand-logo--header {
  flex-direction: row;
  gap: 10px;
}

.ui-brand-logo--auth {
  flex-direction: column;
  gap: 10px;
  text-align: center;
}

.ui-brand-logo__icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: rgb(109 74 255 / 12%);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.ui-brand-logo__icon svg {
  width: 24px;
  height: 24px;
}

.ui-brand-logo__text {
  font-size: 22px;
  font-weight: 700;
  color: var(--tl-color-text);
}

.ui-brand-logo--header .ui-brand-logo__text {
  font-size: 18px;
}
</style>
