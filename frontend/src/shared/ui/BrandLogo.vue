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
      <svg viewBox="0 0 48 48" focusable="false">
        <path
          d="M42.6667 0H5.33333C3.91885 0 2.56229 0.561903 1.5621 1.5621C0.561903 2.56229 0 3.91885 0 5.33333V42.6667C0 44.0812 0.561903 45.4377 1.5621 46.4379C2.56229 47.4381 3.91885 48 5.33333 48H42.6667C44.0812 48 45.4377 47.4381 46.4379 46.4379C47.4381 45.4377 48 44.0812 48 42.6667V5.33333C48 3.91885 47.4381 2.56229 46.4379 1.5621C45.4377 0.561903 44.0812 0 42.6667 0ZM21.3333 34.6667H18.6667C14.3733 34.6667 8 31.84 8 24C8 16.1867 14.3733 13.3333 18.6667 13.3333H21.3333V18.6667H18.6667C17.44 18.6667 13.3333 19.12 13.3333 24C13.3333 29.0667 17.7867 29.3333 18.6667 29.3333H21.3333V34.6667ZM32 21.3333V26.6667H16V21.3333H32ZM29.3333 34.6667H26.6667V29.3333H29.3333C30.56 29.3333 34.6667 28.88 34.6667 24C34.6667 18.9333 30.2133 18.6667 29.3333 18.6667H26.6667V13.3333H29.3333C33.6267 13.3333 40 16.1867 40 24C40 31.84 33.6267 34.6667 29.3333 34.6667Z"
          fill="#6D4AFF"
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.ui-brand-logo__icon svg {
  width: 40px;
  height: 40px;
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
