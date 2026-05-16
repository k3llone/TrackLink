<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from "vue";
import { UiInput } from "@/shared/ui";

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
  change: [q: string];
}>();

const searchValue = ref(props.q);
let timeoutId: ReturnType<typeof setTimeout> | null = null;

const clearPendingChange = () => {
  if (timeoutId) {
    clearTimeout(timeoutId);
    timeoutId = null;
  }
};

const emitChange = () => {
  emit("change", searchValue.value.trim());
};

watch(
  () => props.q,
  (value) => {
    if (searchValue.value !== value) {
      searchValue.value = value;
    }
  },
);

watch(searchValue, () => {
  clearPendingChange();

  if (searchValue.value.trim() === props.q.trim()) {
    return;
  }

  timeoutId = setTimeout(emitChange, 350);
});

onBeforeUnmount(clearPendingChange);
</script>

<template>
  <section class="admin-links-search" aria-label="Поиск ссылок">
    <UiInput
      v-model="searchValue"
      type="search"
      label="Поиск"
      placeholder="UUID, short code или alias"
      autocomplete="off"
      :loading="loading"
    />
  </section>
</template>

<style scoped>
.admin-links-search {
  width: min(100%, 420px);
}

@media (max-width: 767px) {
  .admin-links-search {
    width: 100%;
  }
}
</style>
