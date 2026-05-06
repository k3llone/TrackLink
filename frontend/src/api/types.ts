export type ApiErrorStatus = 400 | 401 | 403 | 404 | 409 | 500;

export interface ApiFieldErrors {
  [field: string]: string;
}

export interface ApiErrorPayload {
  error?: {
    code?: string;
    message?: string;
    fields?: ApiFieldErrors | null;
  };
}

export interface ApiClientError extends Error {
  status: number;
  code: string;
  details: string;
  fields: ApiFieldErrors | null;
}
