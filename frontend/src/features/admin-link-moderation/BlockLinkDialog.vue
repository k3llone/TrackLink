<script setup lang="ts">
import { computed } from "vue";
import type { AdminLink } from "@/api/admin";
import { UiConfirmDialog } from "@/shared/ui";

const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    link: AdminLink | null;
    loading?: boolean;
  }>(),
  {
    loading: false,
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  confirm: [];
}>();

const model = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit("update:modelValue", value),
});

const linkLabel = computed(() => {
  if (!props.link) {
    return "selected link";
  }

  return props.link.shortUrl || props.link.customAlias || props.link.code || props.link.id;
});

const description = computed(() => {
  if (!props.link) {
    return "Selected link will be blocked for redirects.";
  }

  return `Ссылка ${linkLabel.value} будет заблокирована. Redirect на ${props.link.targetUrl} станет недоступен.`;
});
</script>

<template>
  <UiConfirmDialog
    v-model="model"
    title="Заблокировать ссылку?"
    :description="description"
    confirm-text="Заблокировать"
    cancel-text="Отмена"
    variant="danger"
    :loading="loading"
    @confirm="emit('confirm')"
  />
</template>
