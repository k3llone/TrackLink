<script setup lang="ts">
import { computed } from "vue";

type StatusType = "active" | "inactive" | "blocked" | "deleted" | "pending";

const props = defineProps<{
  status: StatusType;
  label?: string;
}>();

const fallbackLabel: Record<StatusType, string> = {
  active: "Active",
  inactive: "Inactive",
  blocked: "Blocked",
  deleted: "Deleted",
  pending: "Pending",
};

const displayLabel = computed(() => props.label || fallbackLabel[props.status]);
</script>

<template>
  <span class="ui-status-badge" :class="`ui-status-badge--${status}`">
    {{ displayLabel }}
  </span>
</template>

<style scoped>
.ui-status-badge {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 700;
}

.ui-status-badge--active {
  color: #25633b;
  background: #d8f1e2;
}

.ui-status-badge--inactive {
  color: #6a6382;
  background: #ebe8f3;
}

.ui-status-badge--blocked {
  color: #7f2631;
  background: #f7d5da;
}

.ui-status-badge--deleted {
  color: #68414d;
  background: #eedfe4;
}

.ui-status-badge--pending {
  color: #8d6025;
  background: #f6e4c2;
}
</style>
