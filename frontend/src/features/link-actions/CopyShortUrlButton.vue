<script setup lang="ts">
import { computed, ref } from "vue";
import { useToast } from "@/shared/composables/useToast";
import { UiButton } from "@/shared/ui";

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
    label: "Копировать",
  },
);

const toast = useToast();
const isCopying = ref(false);

const valueToCopy = computed(() => props.shortUrl.trim());
const isDisabled = computed(() => props.disabled || isCopying.value || !valueToCopy.value);

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
  if (isDisabled.value) {
    return;
  }

  isCopying.value = true;

  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(valueToCopy.value);
    } else {
      fallbackCopy(valueToCopy.value);
    }

    toast.success("Short URL скопирован.");
  } catch {
    toast.error("Не удалось скопировать short URL. Скопируйте его вручную.");
  } finally {
    isCopying.value = false;
  }
};
</script>

<template>
  <UiButton
    type="button"
    :variant="variant"
    :size="size"
    :disabled="isDisabled"
    :loading="isCopying"
    @click.stop="copyShortUrl"
  >
    {{ label }}
  </UiButton>
</template>
