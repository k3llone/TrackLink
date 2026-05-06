import { readonly, ref } from "vue";

type ConfirmVariant = "primary" | "danger";

export type ConfirmOptions = {
  title: string;
  description?: string;
  confirmText?: string;
  cancelText?: string;
  variant?: ConfirmVariant;
};

type ConfirmState = {
  isOpen: boolean;
  loading: boolean;
  options: ConfirmOptions | null;
};

const state = ref<ConfirmState>({
  isOpen: false,
  loading: false,
  options: null,
});

let resolver: ((value: boolean) => void) | null = null;

const open = (options: ConfirmOptions) => {
  state.value = {
    isOpen: true,
    loading: false,
    options,
  };
};

const close = () => {
  state.value = {
    isOpen: false,
    loading: false,
    options: null,
  };
};

const setLoading = (loading: boolean) => {
  state.value = {
    ...state.value,
    loading,
  };
};

const confirm = () => {
  resolver?.(true);
  resolver = null;
  close();
};

const cancel = () => {
  resolver?.(false);
  resolver = null;
  close();
};

const request = (options: ConfirmOptions) =>
  new Promise<boolean>((resolve) => {
    resolver = resolve;
    open(options);
  });

export const useConfirm = () => ({
  state: readonly(state),
  open,
  close,
  confirm,
  cancel,
  request,
  setLoading,
});
