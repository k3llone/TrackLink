<script setup lang="ts">
import type { Link } from "@/entities/link/link.types";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton, UiInput } from "@/shared/ui";
import { useCreateLinkForm } from "./useCreateLinkForm";

const emit = defineEmits<{
  created: [link: Link];
}>();

const { form, errors, isSubmitting, submit } = useCreateLinkForm();
const { t } = useI18n();

const onSubmit = async () => {
  const link = await submit();

  if (link) {
    emit("created", link);
  }
};
</script>

<template>
  <form class="create-link-form" novalidate @submit.prevent="onSubmit">
    <p v-if="errors.form" class="create-link-form__error" role="alert">{{ errors.form }}</p>

    <UiInput
      v-model="form.targetUrl"
      :label="t('common.targetUrl')"
      type="url"
      placeholder="https://example.com/landing"
      autocomplete="url"
      required
      :disabled="isSubmitting"
      :error="errors.targetUrl"
    />

    <UiInput
      v-model="form.customAlias"
      :label="t('createLink.form.aliasLabel')"
      placeholder="spring-campaign"
      :hint="t('createLink.form.aliasHint')"
      autocomplete="off"
      :disabled="isSubmitting"
      :error="errors.customAlias"
    />

    <div class="create-link-form__actions">
      <UiButton type="submit" size="lg" :loading="isSubmitting">{{ t("createLink.form.submit") }}</UiButton>
    </div>
  </form>
</template>

<style scoped>
.create-link-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: min(100%, var(--tl-form-width));
  padding: 22px;
  border: 1px solid rgb(37 31 63 / 10%);
  border-radius: var(--tl-radius-md);
  background: var(--tl-color-white);
}

.create-link-form__error {
  border-radius: var(--tl-radius-md);
  padding: 10px 12px;
  background: rgb(225 75 90 / 10%);
  color: var(--tl-color-danger);
  font-size: 13px;
}

.create-link-form__actions {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

@media (max-width: 640px) {
  .create-link-form {
    padding: 18px;
  }

  .create-link-form__actions,
  .create-link-form__actions :deep(.ui-button) {
    width: 100%;
  }
}
</style>
