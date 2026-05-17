<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import {
  ANALYTICS_PERIOD_OPTIONS,
  DEFAULT_ANALYTICS_PERIOD,
  type AnalyticsPeriodValue,
} from "./analyticsPeriod";

const props = withDefaults(
  defineProps<{
    modelValue?: AnalyticsPeriodValue;
    loading?: boolean;
  }>(),
  {
    modelValue: DEFAULT_ANALYTICS_PERIOD,
    loading: false,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: AnalyticsPeriodValue];
  change: [value: AnalyticsPeriodValue];
}>();

const { t } = useI18n();
const selectedValue = computed(() => props.modelValue);

const optionLabels: Record<AnalyticsPeriodValue, () => string> = {
  "24h": () => t("analytics.period.24h"),
  "7d": () => t("analytics.period.7d"),
  "30d": () => t("analytics.period.30d"),
};

const selectPeriod = (value: AnalyticsPeriodValue) => {
  if (props.loading || value === selectedValue.value) {
    return;
  }

  emit("update:modelValue", value);
  emit("change", value);
};
</script>

<template>
  <section class="analytics-period-picker" :aria-label="t('analytics.period.aria')">
    <span class="analytics-period-picker__label">{{ t("analytics.period.label") }}</span>

    <div class="analytics-period-picker__control" role="group" :aria-label="t('analytics.period.controlAria')">
      <button
        v-for="option in ANALYTICS_PERIOD_OPTIONS"
        :key="option.value"
        class="analytics-period-picker__option"
        :class="{ 'is-selected': option.value === selectedValue }"
        type="button"
        :aria-pressed="option.value === selectedValue"
        :disabled="loading"
        @click="selectPeriod(option.value)"
      >
        {{ optionLabels[option.value]() }}
      </button>
    </div>
  </section>
</template>

<style scoped>
.analytics-period-picker {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.analytics-period-picker__label {
  color: var(--tl-color-text-muted);
  font-size: 13px;
  font-weight: 700;
}

.analytics-period-picker__control {
  display: inline-flex;
  padding: 3px;
  border: 1px solid #ddd7e8;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
}

.analytics-period-picker__option {
  min-height: 34px;
  padding: 8px 12px;
  border: 0;
  border-radius: var(--tl-radius-sm);
  background: transparent;
  color: var(--tl-color-text-muted);
  font-family: var(--tl-font-family);
  font-size: 13px;
  font-weight: 700;
  line-height: 1;
  cursor: pointer;
  transition: background-color 0.2s ease, color 0.2s ease, opacity 0.2s ease;
}

.analytics-period-picker__option:hover:not(:disabled) {
  color: var(--tl-color-text);
}

.analytics-period-picker__option.is-selected {
  background: var(--tl-color-white);
  color: var(--tl-color-primary);
  box-shadow: 0 1px 4px rgb(37 31 63 / 10%);
}

.analytics-period-picker__option:disabled {
  cursor: not-allowed;
  opacity: 0.65;
}

.analytics-period-picker__option:focus-visible {
  outline: 2px solid rgb(109 74 255 / 28%);
  outline-offset: 2px;
}

@media (max-width: 767px) {
  .analytics-period-picker {
    align-items: stretch;
    flex-direction: column;
    width: 100%;
  }

  .analytics-period-picker__control {
    width: 100%;
  }

  .analytics-period-picker__option {
    flex: 1;
  }
}
</style>
