<script setup lang="ts">
import { computed } from "vue";
import type { AdminLink } from "@/api/admin";
import { useI18n } from "@/shared/composables/useI18n";
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

const { t } = useI18n();
const model = computed({
  get: () => props.modelValue,
  set: (value: boolean) => emit("update:modelValue", value),
});

const linkLabel = computed(() => {
  if (!props.link) {
    return t("common.notAvailable");
  }

  return props.link.shortUrl || props.link.customAlias || props.link.code || props.link.id;
});

const description = computed(() => {
  if (!props.link) {
    return t("admin.block.description");
  }

  return t("admin.block.descriptionWithTarget", {
    link: linkLabel.value,
    targetUrl: props.link.targetUrl,
  });
});
</script>

<template>
  <UiConfirmDialog
    v-model="model"
    :title="t('admin.block.confirmTitle')"
    :description="description"
    :confirm-text="t('admin.block.label')"
    :cancel-text="t('common.cancel')"
    variant="danger"
    :loading="loading"
    @confirm="emit('confirm')"
  />
</template>
