<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import { UiPageState } from "@/shared/ui";
import type { LinkUnavailableReason } from "./linkUnavailable.types";

const props = defineProps<{
  reason: LinkUnavailableReason;
}>();

const { t } = useI18n();

const messages: Record<LinkUnavailableReason, { type: "error" | "forbidden" | "not-found"; title: () => string; description: () => string }> = {
  not_found: {
    type: "not-found",
    title: () => t("publicLink.notFound.title"),
    description: () => t("publicLink.notFound.description"),
  },
  inactive: {
    type: "error",
    title: () => t("publicLink.unavailable.title"),
    description: () => t("publicLink.inactive.description"),
  },
  blocked: {
    type: "forbidden",
    title: () => t("publicLink.unavailable.title"),
    description: () => t("publicLink.blocked.description"),
  },
  deleted: {
    type: "error",
    title: () => t("publicLink.unavailable.title"),
    description: () => t("publicLink.deleted.description"),
  },
  gone: {
    type: "error",
    title: () => t("publicLink.unavailable.title"),
    description: () => t("publicLink.gone.description"),
  },
};

const message = computed(() => messages[props.reason]);
</script>

<template>
  <UiPageState :type="message.type" :title="message.title()" :description="message.description()" />
</template>
