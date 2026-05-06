import { readonly, ref } from "vue";
import type { UiToastItem } from "@/shared/ui/Toast.vue";

type ToastInput = Omit<UiToastItem, "id"> & { id?: string };

const DEFAULT_DURATION = 3500;
const toasts = ref<UiToastItem[]>([]);
const timeoutMap = new Map<string, ReturnType<typeof setTimeout>>();

const createId = () => `toast-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;

const removeTimeout = (id: string) => {
  const timeout = timeoutMap.get(id);
  if (!timeout) {
    return;
  }
  clearTimeout(timeout);
  timeoutMap.delete(id);
};

const remove = (id: string) => {
  removeTimeout(id);
  toasts.value = toasts.value.filter((toast) => toast.id !== id);
};

const show = (toast: ToastInput) => {
  const id = toast.id ?? createId();
  const item: UiToastItem = {
    id,
    type: toast.type,
    title: toast.title,
    message: toast.message,
    duration: toast.duration ?? DEFAULT_DURATION,
  };

  toasts.value = [...toasts.value, item];
  removeTimeout(id);

  if (item.duration && item.duration > 0) {
    timeoutMap.set(
      id,
      setTimeout(() => {
        remove(id);
      }, item.duration),
    );
  }

  return id;
};

const clear = () => {
  timeoutMap.forEach((timeout) => clearTimeout(timeout));
  timeoutMap.clear();
  toasts.value = [];
};

const success = (message: string, title?: string) => show({ type: "success", message, title });
const error = (message: string, title?: string) => show({ type: "error", message, title });
const info = (message: string, title?: string) => show({ type: "info", message, title });
const warning = (message: string, title?: string) => show({ type: "warning", message, title });

export const useToast = () => ({
  toasts: readonly(toasts),
  show,
  success,
  error,
  info,
  warning,
  remove,
  clear,
});
