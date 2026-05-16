<script setup lang="ts">
import { watch } from "vue";
import { UiInput } from "@/shared/ui";
import { useLinksSearch, type LinksSearchFilters } from "./useLinksSearch";

const props = withDefaults(
  defineProps<{
    q?: string;
    loading?: boolean;
  }>(),
  {
    q: "",
    loading: false,
  },
);

const emit = defineEmits<{
  change: [filters: LinksSearchFilters];
}>();

const search = useLinksSearch({
  initialQ: props.q,
  onChange: (filters) => emit("change", filters),
});

watch(
  () => props.q,
  (value) => {
    search.setQ(value);
  },
);
</script>

<template>
  <section class="links-search" aria-label="Поиск ссылок">
    <UiInput
      v-model="search.q.value"
      class="links-search__input"
      type="search"
      placeholder="Поиск по short code, alias или target URL"
      autocomplete="off"
      :loading="loading"
    />
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

@media (max-width: 767px) {
  .links-search {
    align-items: stretch;
    flex-direction: column;
    width: 100%;
  }

  .links-search__input {
    min-width: 0;
  }
}
</style>
