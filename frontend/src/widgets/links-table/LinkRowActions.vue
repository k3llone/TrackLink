<script setup lang="ts">
import { useRouter } from "vue-router";
import type { Link } from "@/entities/link/link.types";
import CopyShortUrlButton from "@/features/link-actions/CopyShortUrlButton.vue";
import DeleteLinkButton from "@/features/link-actions/DeleteLinkButton.vue";
import UpdateLinkStatusButton from "@/features/link-actions/UpdateLinkStatusButton.vue";
import { getLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiButton } from "@/shared/ui";

const props = defineProps<{
  link: Link;
}>();

const emit = defineEmits<{
  deleted: [linkId: string];
  updated: [link: Link];
}>();

const router = useRouter();

const openAnalytics = () => {
  void router.push(getLinkDetailsPath(props.link.id));
};

const onUpdated = (link: Link) => {
  emit("updated", link);
};

const onDeleted = (linkId: string) => {
  emit("deleted", linkId);
};
</script>

<template>
  <div class="link-row-actions">
    <CopyShortUrlButton :short-url="link.shortUrl" />
    <UpdateLinkStatusButton :link="link" @updated="onUpdated" />
    <UiButton variant="ghost" size="sm" type="button" @click.stop="openAnalytics">Аналитика</UiButton>
    <DeleteLinkButton :link-id="link.id" :short-url="link.shortUrl" @deleted="onDeleted" />
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
