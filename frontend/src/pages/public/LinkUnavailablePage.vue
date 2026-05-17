<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import LinkUnavailableMessage from "@/widgets/link-unavailable/LinkUnavailableMessage.vue";
import { isLinkUnavailableReason, type LinkUnavailableReason } from "@/widgets/link-unavailable/linkUnavailable.types";

const route = useRoute();

const reason = computed<LinkUnavailableReason>(() => {
  const rawReason = Array.isArray(route.query.reason) ? route.query.reason[0] : route.query.reason;

  // TODO(FR-042): coordinate backend handoff: serve this page for 403/404/410,
  // redirect to /link-unavailable?reason=..., or inject status into this SPA route.
  if (typeof rawReason === "string" && isLinkUnavailableReason(rawReason)) {
    return rawReason;
  }

  return "not_found";
});
</script>

<template>
  <LinkUnavailableMessage :reason="reason" />
</template>
