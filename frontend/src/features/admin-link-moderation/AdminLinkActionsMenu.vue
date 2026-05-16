<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import type { AdminLink } from "@/api/admin";
import { getAdminLinkDetailsPath } from "@/shared/lib/routes/paths";
import { UiButton } from "@/shared/ui";
import BlockLinkButton from "./BlockLinkButton.vue";

const props = defineProps<{
  link: AdminLink;
}>();

const emit = defineEmits<{
  blocked: [link: AdminLink];
}>();

const router = useRouter();
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

const openAnalytics = () => {
  close();
  void router.push(getAdminLinkDetailsPath(props.link.id));
};

const onBlocked = (link: AdminLink) => {
  close();
  emit("blocked", link);
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
      <button
        type="button"
        class="admin-link-actions-menu__item"
        role="menuitem"
        @click.stop="openAnalytics"
      >
        Открыть аналитику
      </button>

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

.admin-link-actions-menu__item {
  width: 100%;
  min-height: 34px;
  border: 0;
  border-radius: var(--tl-radius-md);
  background: transparent;
  color: var(--tl-color-text);
  cursor: pointer;
  font-family: var(--tl-font-family);
  font-size: 13px;
  font-weight: 600;
  text-align: left;
  padding: 8px 12px;
}

.admin-link-actions-menu__item:hover,
.admin-link-actions-menu__item:focus-visible {
  background: var(--tl-color-surface-muted);
  outline: 0;
}

</style>
