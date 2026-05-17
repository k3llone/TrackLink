<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";
import LinkUnavailableMessage from "@/widgets/link-unavailable/LinkUnavailableMessage.vue";
import { normalizeLinkUnavailableReason, type LinkUnavailableReason } from "@/widgets/link-unavailable/linkUnavailable.types";

const route = useRoute();

const reason = computed<LinkUnavailableReason>(() => {
  const rawParamReason = Array.isArray(route.params.reason) ? route.params.reason[0] : route.params.reason;
  const rawQueryReason = Array.isArray(route.query.reason) ? route.query.reason[0] : route.query.reason;
  const rawReason = rawParamReason ?? rawQueryReason;

  return normalizeLinkUnavailableReason(rawReason);
});
</script>

<template>
  <LinkUnavailableMessage :reason="reason" />
</template>
