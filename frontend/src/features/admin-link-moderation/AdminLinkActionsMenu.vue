<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from "vue";
import type { AdminLink } from "@/api/admin";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton } from "@/shared/ui";
import BlockLinkButton from "./BlockLinkButton.vue";
import UnblockLinkButton from "./UnblockLinkButton.vue";

defineProps<{
  link: AdminLink;
}>();

const emit = defineEmits<{
  blocked: [link: AdminLink];
  unblocked: [link: AdminLink];
}>();

const root = ref<HTMLElement | null>(null);
const menu = ref<HTMLElement | null>(null);
const isOpen = ref(false);
const menuStyle = ref<Record<string, string>>({});
const { t } = useI18n();

const close = () => {
  isOpen.value = false;
};

const updateMenuPosition = () => {
  if (!isOpen.value || !root.value) {
    return;
  }

  const rect = root.value.getBoundingClientRect();
  const gap = 6;
  const viewportMargin = 12;
  const menuWidth = menu.value?.offsetWidth ?? 190;
  const menuHeight = menu.value?.offsetHeight ?? 90;
  const left = Math.max(
    viewportMargin,
    Math.min(window.innerWidth - menuWidth - viewportMargin, rect.right - menuWidth),
  );
  const topBelow = rect.bottom + gap;
  const top =
    topBelow + menuHeight > window.innerHeight - viewportMargin
      ? Math.max(viewportMargin, rect.top - menuHeight - gap)
      : topBelow;

  menuStyle.value = {
    left: `${Math.round(left)}px`,
    top: `${Math.round(top)}px`,
    visibility: "visible",
  };
};

const open = async () => {
  menuStyle.value = {
    left: "0",
    top: "0",
    visibility: "hidden",
  };
  isOpen.value = true;
  await nextTick();
  updateMenuPosition();
};

const toggle = () => {
  if (isOpen.value) {
    close();
    return;
  }

  void open();
};

const onDocumentClick = (event: MouseEvent) => {
  const target = event.target as Node;

  if (target instanceof Element && target.closest(".ui-modal")) {
    return;
  }

  if (!root.value?.contains(target) && !menu.value?.contains(target)) {
    close();
  }
};

const onDocumentKeydown = (event: KeyboardEvent) => {
  if (event.key === "Escape") {
    close();
  }
};

const onViewportChange = () => {
  updateMenuPosition();
};

const onBlocked = (link: AdminLink) => {
  close();
  emit("blocked", link);
};

const onUnblocked = (link: AdminLink) => {
  close();
  emit("unblocked", link);
};

onMounted(() => {
  document.addEventListener("click", onDocumentClick);
  document.addEventListener("keydown", onDocumentKeydown);
  document.addEventListener("scroll", onViewportChange, true);
  window.addEventListener("resize", onViewportChange);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", onDocumentClick);
  document.removeEventListener("keydown", onDocumentKeydown);
  document.removeEventListener("scroll", onViewportChange, true);
  window.removeEventListener("resize", onViewportChange);
});
</script>

<template>
  <div ref="root" class="admin-link-actions-menu" @click.stop>
    <UiButton
      type="button"
      variant="ghost"
      size="sm"
      class="admin-link-actions-menu__trigger"
      :aria-label="t('admin.actions.open')"
      :aria-expanded="isOpen"
      @click.stop="toggle"
    >
      ...
    </UiButton>

    <Teleport to="body">
      <div
        v-if="isOpen"
        ref="menu"
        class="admin-link-actions-menu__menu"
        role="menu"
        :style="menuStyle"
        @click.stop
      >
        <UnblockLinkButton
          v-if="link.status === 'blocked'"
          :link="link"
          variant="ghost"
          size="sm"
          full-width
          @unblocked="onUnblocked"
        />

        <BlockLinkButton
          v-else
          :link="link"
          variant="ghost"
          size="sm"
          full-width
          @blocked="onBlocked"
        />
      </div>
    </Teleport>
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
  position: fixed;
  z-index: 1200;
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
