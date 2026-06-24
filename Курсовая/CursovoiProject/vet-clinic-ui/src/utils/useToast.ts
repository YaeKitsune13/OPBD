import { ref } from "vue";

export interface Toast {
  id: number;
  message: string;
  type: "success" | "error" | "info";
}

const toasts = ref<Toast[]>([]);
let counter = 0;

export function useToast() {
  const showToast = (
    message: string,
    type: "success" | "error" | "info" = "success",
  ) => {
    const id = ++counter;
    const toast: Toast = { id, message, type };

    toasts.value.push(toast);

    setTimeout(() => {
      removeToast(id);
    }, 3000);
  };

  const removeToast = (id: number) => {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  };

  return { toasts, showToast, removeToast };
}
