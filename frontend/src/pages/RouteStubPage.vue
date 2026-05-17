<script setup lang="ts">
import { useI18n } from "@/shared/composables/useI18n";

defineProps<{
  title: string;
}>();

const { t } = useI18n();
</script>

<template>
  <section>
    <h1>{{ title }}</h1>
    <p class="text-muted">{{ t("routeStub.description") }}</p>
  </section>
</template>
