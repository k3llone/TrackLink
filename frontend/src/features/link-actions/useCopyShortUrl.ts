import { ref } from "vue";
import { useI18n } from "@/shared/composables/useI18n";
import { useToast } from "@/shared/composables/useToast";

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

export const useCopyShortUrl = () => {
  const toast = useToast();
  const { t } = useI18n();
  const isCopying = ref(false);

  const copyShortUrl = async (shortUrl: string) => {
    const valueToCopy = shortUrl.trim();

    if (isCopying.value || !valueToCopy) {
      return;
    }

    isCopying.value = true;

    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(valueToCopy);
      } else {
        fallbackCopy(valueToCopy);
      }

      toast.success(t("linkActions.copy.success"));
    } catch {
      toast.error(t("linkActions.copy.error"));
    } finally {
      isCopying.value = false;
    }
  };

  return {
    isCopying,
    copyShortUrl,
  };
};
