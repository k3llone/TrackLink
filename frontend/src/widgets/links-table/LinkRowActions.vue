<script setup lang="ts">
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { deleteLink } from "@/api/links";
import type { ApiClientError } from "@/api/types";
import type { Link } from "@/entities/link/link.types";
import { useToast } from "@/shared/composables/useToast";
import { getLinkEditPath } from "@/shared/lib/routes/paths";
import { UiButton, UiConfirmDialog } from "@/shared/ui";

const props = defineProps<{
  link: Link;
}>();

const emit = defineEmits<{
  deleted: [link: Link];
}>();

const router = useRouter();
const toast = useToast();
const isConfirmOpen = ref(false);
const isDeleting = ref(false);

const isDeleted = computed(() => props.link.status === "deleted");
const isBlocked = computed(() => props.link.status === "blocked");
const canEdit = computed(() => !isDeleted.value && !isBlocked.value);
const canExport = computed(() => !isDeleted.value);
const canDelete = computed(() => !isDeleted.value);

const unavailableByStatusText = "Действие недоступно для этого статуса.";
const editTitle = computed(() => (canEdit.value ? "Edit link" : unavailableByStatusText));
const exportTitle = computed(() => (canExport.value ? "Export link data" : unavailableByStatusText));
const deleteTitle = computed(() => (canDelete.value ? "Delete link" : unavailableByStatusText));

const escapeCsvValue = (value: string | number | null | undefined) => `"${String(value ?? "").replaceAll('"', '""')}"`;

const editLink = () => {
  if (canEdit.value) {
    void router.push(getLinkEditPath(props.link.id));
  }
};

const exportLink = () => {
  if (!canExport.value) {
    return;
  }

  const headers = ["id", "shortUrl", "targetUrl", "createdAt", "status", "totalClicks"];
  const values = [
    props.link.id,
    props.link.shortUrl || props.link.code,
    props.link.targetUrl,
    props.link.createdAt,
    props.link.status,
    props.link.totalClicks,
  ];
  const csv = `${headers.join(",")}\n${values.map(escapeCsvValue).join(",")}\n`;
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = `tracklink-${props.link.code}.csv`;
  anchor.click();
  URL.revokeObjectURL(url);

  toast.success("Данные ссылки экспортированы.");
};

const requestDelete = () => {
  if (canDelete.value) {
    isConfirmOpen.value = true;
  }
};

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

const getDeleteErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 404) {
      return "Ссылка уже удалена или недоступна.";
    }

    if (error.status === 401) {
      return "Сессия недействительна. Войдите заново и повторите удаление.";
    }
  }

  return "Не удалось удалить ссылку. Повторите попытку позже.";
};

const confirmDelete = async () => {
  if (isDeleting.value) {
    return;
  }

  isDeleting.value = true;

  try {
    await deleteLink(props.link.id);
    isConfirmOpen.value = false;
    toast.success("Ссылка удалена.");
    emit("deleted", props.link);
  } catch (error: unknown) {
    toast.error(getDeleteErrorMessage(error));
  } finally {
    isDeleting.value = false;
  }
};
</script>

<template>
  <div class="link-row-actions" @click.stop>
    <UiButton variant="primary" size="sm" type="button" :disabled="!canEdit" :title="editTitle" @click="editLink">
      Edit
    </UiButton>
    <UiButton variant="primary" size="sm" type="button" :disabled="!canExport" :title="exportTitle" @click="exportLink">
      Export
    </UiButton>
    <UiButton
      variant="danger"
      size="sm"
      type="button"
      :disabled="!canDelete"
      :loading="isDeleting"
      :title="deleteTitle"
      @click="requestDelete"
    >
      Delete
    </UiButton>

    <UiConfirmDialog
      v-model="isConfirmOpen"
      title="Удалить ссылку?"
      description="Ссылка будет удалена из списка. Это действие нельзя отменить."
      confirm-text="Delete"
      cancel-text="Cancel"
      :loading="isDeleting"
      @confirm="confirmDelete"
    />
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
