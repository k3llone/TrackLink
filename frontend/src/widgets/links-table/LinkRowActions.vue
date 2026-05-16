<script setup lang="ts">
import { useRouter } from "vue-router";
import type { Link } from "@/entities/link/link.types";
import UpdateLinkStatusButton from "@/features/link-actions/UpdateLinkStatusButton.vue";
import { getLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiButton } from "@/shared/ui";

const props = defineProps<{
  link: Link;
}>();

const emit = defineEmits<{
  updated: [link: Link];
}>();

const router = useRouter();

const openAnalytics = () => {
  void router.push(getLinkDetailsPath(props.link.id));
};

const onUpdated = (link: Link) => {
  emit("updated", link);
};
</script>

<template>
  <div class="link-row-actions">
    <UpdateLinkStatusButton :link="link" @updated="onUpdated" />
    <UiButton variant="ghost" size="sm" type="button" @click.stop="openAnalytics">Аналитика</UiButton>
  </div>
</template>

<style scoped>
.link-row-actions {
  display: inline-flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  flex-wrap: wrap;
}
</style>
