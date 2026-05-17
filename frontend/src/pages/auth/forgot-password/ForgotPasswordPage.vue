<script setup lang="ts">
import { useI18n } from "@/shared/composables/useI18n";
import { UiPageState } from "@/shared/ui";

const { t } = useI18n();
</script>

<template>
  <UiPageState
    type="empty"
    :title="t('auth.forgot.title')"
    :description="t('auth.forgot.description')"
  />
</template>
