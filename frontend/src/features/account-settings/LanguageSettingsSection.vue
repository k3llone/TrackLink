<script setup lang="ts">
import { computed } from "vue";
import { useToast } from "@/shared/composables/useToast";
import { useI18n } from "@/shared/composables/useI18n";
import type { Locale } from "@/shared/lib/i18n";

const { locale, setLocale, t } = useI18n();
const toast = useToast();

const languageOptions = computed<Array<{ value: Locale; label: string; description: string }>>(() => [
  {
    value: "en",
    label: t("common.language.english"),
    description: t("settings.language.englishDescription"),
  },
  {
    value: "ru",
    label: t("common.language.russian"),
    description: t("settings.language.russianDescription"),
  },
]);

const selectLocale = (nextLocale: Locale) => {
  if (nextLocale === locale.value) {
    return;
  }

  setLocale(nextLocale);
  const languageKey = nextLocale === "ru" ? "common.language.russian" : "common.language.english";
  toast.success(t("settings.language.changedMessage", { language: t(languageKey) }), t("settings.language.changedTitle"));
};
</script>

<template>
  <section class="language-settings" aria-labelledby="language-settings-title">
    <div class="language-settings__copy">
      <h2 id="language-settings-title" class="language-settings__title">{{ t("settings.language.title") }}</h2>
      <p class="language-settings__subtitle">{{ t("settings.language.subtitle") }}</p>
    </div>

    <div class="language-settings__options" role="group" :aria-label="t('settings.language.ariaLabel')">
      <button
        v-for="option in languageOptions"
        :key="option.value"
        type="button"
        class="language-settings__option"
        :class="{ 'is-active': option.value === locale }"
        :aria-pressed="option.value === locale"
        @click="selectLocale(option.value)"
      >
        <span class="language-settings__option-label">{{ option.label }}</span>
        <span class="language-settings__option-description">{{ option.description }}</span>
      </button>
    </div>
  </section>
</template>

<style scoped>
.language-settings {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
  width: min(100%, 950px);
}

.language-settings__copy {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.language-settings__title {
  color: var(--tl-color-text);
  font-size: 16px;
  line-height: 1.3;
}

.language-settings__subtitle {
  color: var(--tl-color-text-muted);
  font-size: 14px;
}

.language-settings__options {
  display: inline-grid;
  grid-template-columns: repeat(2, minmax(140px, 1fr));
  gap: 8px;
  width: min(100%, 380px);
}

.language-settings__option {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-height: 76px;
  padding: 12px;
  border: 1px solid #ddd7e8;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
  color: var(--tl-color-text);
  cursor: pointer;
  font-family: var(--tl-font-family);
  text-align: left;
  transition: background-color 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.language-settings__option:hover {
  background: #eee9f7;
}

.language-settings__option.is-active {
  border-color: var(--tl-color-primary);
  background: var(--tl-color-white);
  box-shadow: 0 0 0 2px rgb(109 74 255 / 14%);
}

.language-settings__option-label {
  font-size: 14px;
  font-weight: 800;
}

.language-settings__option-description {
  color: var(--tl-color-text-muted);
  font-size: 12px;
  line-height: 1.35;
}

@media (max-width: 767px) {
  .language-settings {
    flex-direction: column;
  }

  .language-settings__options {
    grid-template-columns: 1fr;
    width: 100%;
  }
}
</style>
