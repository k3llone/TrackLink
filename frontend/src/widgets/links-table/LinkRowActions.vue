<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import type { Link } from "@/entities/link/link.types";
import { useToast } from "@/shared/composables/useToast";
import { getLinkDetailsPath, getLinkEditPath } from "@/shared/lib/routes/paths";
import { UiButton } from "@/shared/ui";

const props = defineProps<{
  link: Link;
}>();

const router = useRouter();
const toast = useToast();

const isDeleted = computed(() => props.link.status === "deleted");
const isBlocked = computed(() => props.link.status === "blocked");
const canCopy = computed(() => !isDeleted.value);
const canOpenAnalytics = computed(() => !isDeleted.value);
const canEdit = computed(() => !isDeleted.value && !isBlocked.value);
const statusActionLabel = computed(() => (props.link.status === "active" ? "Деактивировать" : "Активировать"));

const unavailableByStatusText = "Действие недоступно для этого статуса.";
const statusPlaceholderText = "Будет доступно после FR-028.";
const deletePlaceholderText = "Будет доступно после FR-029.";

const copyTitle = computed(() => (canCopy.value ? "Скопировать short URL" : unavailableByStatusText));
const analyticsTitle = computed(() => (canOpenAnalytics.value ? "Открыть аналитику" : unavailableByStatusText));
const editTitle = computed(() => (canEdit.value ? "Редактировать ссылку" : unavailableByStatusText));
const statusTitle = computed(() => (isDeleted.value || isBlocked.value ? unavailableByStatusText : statusPlaceholderText));
const deleteTitle = computed(() => (isDeleted.value ? unavailableByStatusText : deletePlaceholderText));

const fallbackCopy = (value: string) => {
  const textarea = document.createElement("textarea");
  textarea.value = value;
  textarea.setAttribute("readonly", "");
  textarea.style.position = "fixed";
  textarea.style.top = "0";
  textarea.style.left = "0";
  textarea.style.opacity = "0";

  document.body.appendChild(textarea);
  textarea.focus();
  textarea.select();

  const isCopied = document.execCommand("copy");
  document.body.removeChild(textarea);

  if (!isCopied) {
    throw new Error("Copy command failed");
  }
};

const copyShortUrl = async () => {
  if (!canCopy.value) {
    return;
  }

  const shortUrl = props.link.shortUrl || props.link.code;

  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(shortUrl);
    } else {
      fallbackCopy(shortUrl);
    }

    toast.success("Short URL скопирован.");
  } catch {
    toast.error("Не удалось скопировать short URL. Скопируйте его вручную.");
  }
};

const openAnalytics = () => {
  if (canOpenAnalytics.value) {
    void router.push(getLinkDetailsPath(props.link.id));
  }
};

const editLink = () => {
  if (canEdit.value) {
    void router.push(getLinkEditPath(props.link.id));
  }
};
</script>

<template>
  <div class="link-row-actions">
    <UiButton variant="ghost" size="sm" type="button" :disabled="!canCopy" :title="copyTitle" @click="copyShortUrl">
      Копировать
    </UiButton>
    <UiButton
      variant="ghost"
      size="sm"
      type="button"
      :disabled="!canOpenAnalytics"
      :title="analyticsTitle"
      @click="openAnalytics"
    >
      Аналитика
    </UiButton>
    <UiButton variant="ghost" size="sm" type="button" :disabled="!canEdit" :title="editTitle" @click="editLink">
      Изменить
    </UiButton>
    <UiButton variant="ghost" size="sm" type="button" disabled :title="statusTitle">
      {{ statusActionLabel }}
    </UiButton>
    <UiButton variant="danger" size="sm" type="button" disabled :title="deleteTitle">
      Удалить
    </UiButton>
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
