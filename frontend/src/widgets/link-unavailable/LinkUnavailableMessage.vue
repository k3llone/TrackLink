<script setup lang="ts">
import { computed } from "vue";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiPageState } from "@/shared/ui";
import type { LinkUnavailableReason } from "./linkUnavailable.types";

const props = defineProps<{
  reason: LinkUnavailableReason;
}>();

const messages: Record<LinkUnavailableReason, { type: "error" | "forbidden" | "not-found"; title: string; description: string }> = {
  not_found: {
    type: "not-found",
    title: "Link not found",
    description: "This short link does not exist or may have been typed incorrectly.",
  },
  inactive: {
    type: "error",
    title: "Link is unavailable",
    description: "This short link is no longer active and cannot be opened.",
  },
  blocked: {
    type: "forbidden",
    title: "Link is unavailable",
    description: "This short link cannot be opened. Please use a different link or contact the sender.",
  },
  deleted: {
    type: "error",
    title: "Link is unavailable",
    description: "This short link is no longer available.",
  },
  gone: {
    type: "error",
    title: "Link is unavailable",
    description: "This short link is no longer available.",
  },
};

const message = computed(() => messages[props.reason]);
</script>

<template>
  <UiPageState
    :type="message.type"
    :title="message.title"
    :description="message.description"
    action-text="Go to sign in"
    :action-to="ROUTES.login"
  />
</template>
