<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import { UiInput } from "@/shared/ui";

const props = withDefaults(
  defineProps<{
    q?: string;
  }>(),
  {
    q: "",
  },
);

const emit = defineEmits<{
  change: [q: string];
}>();

const { t } = useI18n();
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
  <section class="admin-links-search" :aria-label="t('admin.search.aria')">
    <UiInput
      v-model="searchValue"
      type="search"
      :label="t('common.search')"
      :placeholder="t('admin.search.placeholder')"
      autocomplete="off"
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
