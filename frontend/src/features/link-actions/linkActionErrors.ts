import type { ApiClientError } from "@/api/types";
import { t } from "@/shared/lib/i18n";

const isApiClientError = (error: unknown): error is ApiClientError =>
  error instanceof Error && error.name === "ApiClientError" && "status" in error;

export const getUpdateStatusErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 403) {
      return t("linkActions.status.error.forbidden");
    }

    if (error.status === 404) {
      return t("linkActions.status.error.notFound");
    }

    if (error.status === 409) {
      return t("linkActions.status.error.conflict");
    }

    if (error.status >= 500) {
      return t("linkActions.status.error.server");
    }
  }

  return t("linkActions.status.error.failed");
};

export const getDeleteLinkErrorMessage = (error: unknown) => {
  if (isApiClientError(error)) {
    if (error.status === 403) {
      return t("linkActions.delete.error.forbidden");
    }

    if (error.status === 404) {
      return t("linkActions.delete.error.notFound");
    }

    if (error.status === 409) {
      return t("linkActions.delete.error.conflict");
    }

    if (error.status >= 500) {
      return t("linkActions.delete.error.server");
    }
  }

  return t("linkActions.delete.error.failed");
};
