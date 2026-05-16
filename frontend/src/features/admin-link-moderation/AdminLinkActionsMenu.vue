<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";
import type { AdminLink } from "@/api/admin";
import { UiButton } from "@/shared/ui";
import BlockLinkButton from "./BlockLinkButton.vue";
import DeactivateLinkButton from "./DeactivateLinkButton.vue";

defineProps<{
  link: AdminLink;
}>();

const emit = defineEmits<{
  blocked: [link: AdminLink];
  deactivated: [link: AdminLink];
}>();

const root = ref<HTMLElement | null>(null);
const isOpen = ref(false);

const close = () => {
  isOpen.value = false;
};

const toggle = () => {
  isOpen.value = !isOpen.value;
};

const onDocumentClick = (event: MouseEvent) => {
  if (!root.value?.contains(event.target as Node)) {
    close();
  }
};

const onDocumentKeydown = (event: KeyboardEvent) => {
  if (event.key === "Escape") {
    close();
  }
};

const onBlocked = (link: AdminLink) => {
  close();
  emit("blocked", link);
};

const onDeactivated = (link: AdminLink) => {
  close();
  emit("deactivated", link);
};

onMounted(() => {
  document.addEventListener("click", onDocumentClick);
  document.addEventListener("keydown", onDocumentKeydown);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", onDocumentClick);
  document.removeEventListener("keydown", onDocumentKeydown);
});
</script>

<template>
  <div ref="root" class="admin-link-actions-menu" @click.stop>
    <UiButton
      type="button"
      variant="ghost"
      size="sm"
      class="admin-link-actions-menu__trigger"
      aria-label="Открыть действия со ссылкой"
      :aria-expanded="isOpen"
      @click.stop="toggle"
    >
      ...
    </UiButton>

    <div v-if="isOpen" class="admin-link-actions-menu__menu" role="menu">
      <DeactivateLinkButton
        :link="link"
        variant="ghost"
        size="sm"
        full-width
        @deactivated="onDeactivated"
      />

      <BlockLinkButton
        :link="link"
        variant="ghost"
        size="sm"
        full-width
        @blocked="onBlocked"
      />
    </div>
  </div>
</template>

<style scoped>
.admin-link-actions-menu {
  position: relative;
  display: inline-flex;
  justify-content: flex-end;
}

.admin-link-actions-menu__trigger {
  min-width: 34px;
}

.admin-link-actions-menu__menu {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  z-index: 20;
  display: flex;
  min-width: 190px;
  flex-direction: column;
  gap: 4px;
  padding: 6px;
  border: 1px solid #e4dfef;
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-white);
  box-shadow: 0 12px 28px rgb(20 16 38 / 16%);
}

</style>
