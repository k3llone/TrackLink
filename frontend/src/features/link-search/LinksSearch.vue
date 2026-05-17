<script setup lang="ts">
import { computed, watch } from "vue";
import type { LinkStatus } from "@/entities/link/link.types";
import { useI18n } from "@/shared/composables/useI18n";
import { UiInput } from "@/shared/ui";
import { useLinksSearch, type LinksSearchFilters, type LinksSearchStatus } from "./useLinksSearch";

const props = withDefaults(
  defineProps<{
    q?: string;
    status?: LinksSearchStatus;
    loading?: boolean;
  }>(),
  {
    q: "",
    status: "",
    loading: false,
  },
);

const emit = defineEmits<{
  change: [filters: LinksSearchFilters];
}>();

const { t } = useI18n();
const search = useLinksSearch({
  initialQ: props.q,
  initialStatus: props.status,
  onChange: (filters) => emit("change", filters),
});

const statusOptions = computed<Array<{ value: LinksSearchStatus; label: string }>>(() => [
  { value: "", label: t("links.search.allStatuses") },
  { value: "active", label: t("links.search.active") },
  { value: "inactive", label: t("links.search.inactive") },
  { value: "blocked", label: t("links.search.blocked") },
  { value: "deleted", label: t("links.search.deleted") },
]);

watch(
  () => props.q,
  (value) => {
    search.setQ(value);
  },
);

watch(
  () => props.status,
  (value) => {
    search.setStatus(value);
  },
);

const onStatusChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  search.setStatus(target.value as LinkStatus | "");
};
</script>

<template>
  <section class="links-search" :aria-label="t('links.search.aria')">
    <UiInput
      v-model="search.q.value"
      class="links-search__input"
      type="search"
      :placeholder="t('links.search.placeholder')"
      autocomplete="off"
      :loading="loading"
    />

    <label class="links-search__status">
      <span class="links-search__status-label">{{ t("common.status") }}</span>
      <select
        class="links-search__select"
        :value="search.status.value"
        :disabled="loading"
        @change="onStatusChange"
      >
        <option v-for="option in statusOptions" :key="option.value || 'all'" :value="option.value">
          {{ option.label }}
        </option>
      </select>
    </label>
  </section>
</template>

<style scoped>
.links-search {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  width: min(100%, 540px);
}

.links-search__input {
  min-width: 260px;
}

.links-search__status {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 190px;
  color: var(--tl-color-text);
}

.links-search__status-label {
  font-size: 13px;
  font-weight: 700;
}

.links-search__select {
  min-height: 44px;
  padding: 0 12px;
  border: 1px solid #ddd7e8;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-surface-muted);
  color: var(--tl-color-text);
  font-family: var(--tl-font-family);
  font-size: 14px;
}

.links-search__select:disabled {
  opacity: 0.65;
}

.links-search__select:focus {
  border-color: var(--tl-color-primary);
  box-shadow: 0 0 0 2px rgb(109 74 255 / 18%);
  outline: 0;
}

@media (max-width: 767px) {
  .links-search {
    align-items: stretch;
    flex-direction: column;
    width: 100%;
  }

  .links-search__input {
    min-width: 0;
  }

  .links-search__status {
    min-width: 0;
  }
}
</style>
