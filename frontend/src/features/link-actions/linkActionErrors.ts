import type { ApiClientError } from "@/api/types";

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

export const getUpdateStatusErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 403) {
      return "У вас нет доступа к изменению статуса этой ссылки.";
    }

    if (error.status === 404) {
      return "Ссылка не найдена или уже недоступна.";
    }

    if (error.status === 409) {
      return "Этот статус нельзя изменить. Заблокированные и удаленные ссылки недоступны для активации.";
    }

    if (error.status >= 500) {
      return "Сервер временно недоступен. Повторите попытку позже.";
    }
  }

  return "Не удалось изменить статус ссылки. Проверьте соединение и повторите попытку.";
};
