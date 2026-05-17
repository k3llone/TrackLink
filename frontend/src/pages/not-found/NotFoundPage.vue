<script setup lang="ts">
import { useI18n } from "@/shared/composables/useI18n";
import { UiPageState } from "@/shared/ui";

const { t } = useI18n();
</script>

<template>
  <UiPageState
    type="not-found"
    :title="t('notFound.title')"
    :description="t('notFound.description')"
  />
</template>
