import type { ApiClientError, ApiErrorPayload } from "@/api/types";
import { env } from "@/shared/config/env";
import { t, type MessageKey } from "@/shared/lib/i18n";

const DEFAULT_ERROR_KEYS: Record<number, MessageKey> = {
  400: "api.error.400",
  401: "api.error.401",
  403: "api.error.403",
  404: "api.error.404",
  409: "api.error.409",
  500: "api.error.500",
};

type RequestOptions = Omit<RequestInit, "body"> & {
  body?: unknown;
};

function getFallbackMessage(status: number): string {
  const key = DEFAULT_ERROR_KEYS[status] ?? "api.error.default";
  return t(key);
}

async function parseJson<T>(response: Response): Promise<T | null> {
  const contentType = response.headers.get("content-type");
  if (!contentType?.includes("application/json")) {
    return null;
  }

  return (await response.json()) as T;
}

function createApiError(status: number, payload: ApiErrorPayload | null): ApiClientError {
  const code = payload?.error?.code || "api_error";
  const details = payload?.error?.message || getFallbackMessage(status);
  const error = new Error(details) as ApiClientError;

  error.name = "ApiClientError";
  error.status = status;
  error.code = code;
  error.details = details;
  error.fields = payload?.error?.fields || null;

  return error;
}

async function request<TResponse>(path: string, options: RequestOptions = {}): Promise<TResponse> {
  const { body, headers, ...restOptions } = options;
  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    method: restOptions.method || "GET",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
    body: body === undefined ? undefined : JSON.stringify(body),
    ...restOptions,
  });

  if (!response.ok) {
    const payload = await parseJson<ApiErrorPayload>(response);
    throw createApiError(response.status, payload);
  }

  if (response.status === 204) {
    return undefined as TResponse;
  }

  const data = await parseJson<TResponse>(response);
  return (data ?? (undefined as TResponse)) as TResponse;
}

export const http = {
  request,
  get<TResponse>(path: string, options?: Omit<RequestOptions, "method">) {
    return request<TResponse>(path, { ...options, method: "GET" });
  },
  post<TResponse>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">) {
    return request<TResponse>(path, { ...options, method: "POST", body });
  },
  put<TResponse>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">) {
    return request<TResponse>(path, { ...options, method: "PUT", body });
  },
  patch<TResponse>(path: string, body?: unknown, options?: Omit<RequestOptions, "method" | "body">) {
    return request<TResponse>(path, { ...options, method: "PATCH", body });
  },
  delete<TResponse>(path: string, options?: Omit<RequestOptions, "method">) {
    return request<TResponse>(path, { ...options, method: "DELETE" });
  },
};
