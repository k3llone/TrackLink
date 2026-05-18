<script setup lang="ts">
import { computed } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import { UiButton } from "@/shared/ui";
import { useCopyShortUrl } from "./useCopyShortUrl";

type ButtonVariant = "primary" | "secondary" | "danger" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

const props = withDefaults(
  defineProps<{
    shortUrl: string;
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
    label?: string;
  }>(),
  {
    variant: "ghost",
    size: "sm",
    disabled: false,
    label: "",
  },
);

const { t } = useI18n();
const { isCopying, copyShortUrl } = useCopyShortUrl();

const valueToCopy = computed(() => props.shortUrl.trim());
const isDisabled = computed(() => props.disabled || isCopying.value || !valueToCopy.value);
const buttonLabel = computed(() => props.label || t("common.copy"));

const onCopyShortUrl = () => {
  if (isDisabled.value) {
    return;
  }

  void copyShortUrl(valueToCopy.value);
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isCopying"
    @click.stop="onCopyShortUrl"
  >
    {{ buttonLabel }}
  </UiButton>
</template>
