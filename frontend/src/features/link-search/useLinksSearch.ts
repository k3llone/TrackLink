import { onBeforeUnmount, ref, watch } from "vue";

export interface LinksSearchFilters {
  q: string;
}

interface UseLinksSearchOptions {
  initialQ?: string;
  debounceMs?: number;
  onChange?: (filters: LinksSearchFilters) => void;
}

export const useLinksSearch = (options: UseLinksSearchOptions = {}) => {
  const q = ref(options.initialQ ?? "");
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

  onBeforeUnmount(clearPendingChange);

  return {
    q,
    setQ,
  };
};
