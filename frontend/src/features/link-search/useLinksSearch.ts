import { onBeforeUnmount, ref, watch } from "vue";
import type { LinkStatus } from "@/entities/link/link.types";

export type LinksSearchStatus = LinkStatus | "";

export interface LinksSearchFilters {
  q: string;
  status: LinksSearchStatus;
}

interface UseLinksSearchOptions {
  initialQ?: string;
  initialStatus?: LinksSearchStatus;
  debounceMs?: number;
  onChange?: (filters: LinksSearchFilters) => void;
}

export const useLinksSearch = (options: UseLinksSearchOptions = {}) => {
  const q = ref(options.initialQ ?? "");
  const status = ref<LinksSearchStatus>(options.initialStatus ?? "");
  const debounceMs = options.debounceMs ?? 350;
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  const clearPendingChange = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
  };

  const emitChange = () => {
    options.onChange?.({
      q: q.value.trim(),
      status: status.value,
    });
  };

  watch(q, () => {
    clearPendingChange();
    timeoutId = setTimeout(emitChange, debounceMs);
  });

  const setQ = (value: string) => {
    if (q.value !== value) {
      q.value = value;
    }
  };

  const setStatus = (value: LinksSearchStatus) => {
    if (status.value !== value) {
      status.value = value;
    }
  };

  watch(status, () => {
    clearPendingChange();
    emitChange();
  });

  onBeforeUnmount(clearPendingChange);

  return {
    q,
    status,
    setQ,
    setStatus,
  };
};
