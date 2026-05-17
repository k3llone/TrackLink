import { readonly, ref } from "vue";
import { messages } from "./messages";

export const LOCALES = ["en", "ru"] as const;
export const I18N_STORAGE_KEY = "tracklink.locale";

export type Locale = (typeof LOCALES)[number];
export type MessageKey = keyof typeof messages.en;
export type MessageParams = Record<string, string | number>;

const locale = ref<Locale>("en");

const isLocale = (value: unknown): value is Locale =>
  typeof value === "string" && LOCALES.includes(value as Locale);

const getBrowserLocale = (): Locale => {
  if (typeof navigator === "undefined") {
    return "en";
  }

  return navigator.language.toLowerCase().startsWith("ru") ? "ru" : "en";
};

const getStoredLocale = (): Locale | null => {
  if (typeof localStorage === "undefined") {
    return null;
  }

  try {
    const storedLocale = localStorage.getItem(I18N_STORAGE_KEY);
    return isLocale(storedLocale) ? storedLocale : null;
  } catch {
    return null;
  }
};

const persistLocale = (nextLocale: Locale) => {
  if (typeof localStorage === "undefined") {
    return;
  }

  try {
    localStorage.setItem(I18N_STORAGE_KEY, nextLocale);
  } catch {
    // Language preference is non-critical; keep the UI reactive even if storage is unavailable.
  }
};

const applyDocumentLocale = (nextLocale: Locale) => {
  if (typeof document !== "undefined") {
    document.documentElement.lang = nextLocale;
  }
};

const resolveInitialLocale = () => getStoredLocale() ?? getBrowserLocale();

const interpolate = (template: string, params: MessageParams = {}) =>
  template.replace(/\{(\w+)\}/g, (_, key: string) => String(params[key] ?? `{${key}}`));

export const setLocale = (nextLocale: Locale) => {
  locale.value = nextLocale;
  persistLocale(nextLocale);
  applyDocumentLocale(nextLocale);
};

export const getLocale = () => locale.value;

export const t = (key: MessageKey, params?: MessageParams) => {
  const template = messages[locale.value][key] ?? messages.en[key] ?? key;
  return interpolate(template, params);
};

export const formatNumber = (value: number) =>
  new Intl.NumberFormat(locale.value === "ru" ? "ru-RU" : "en-US").format(value);

export const formatDate = (
  value: string | number | Date,
  options: Intl.DateTimeFormatOptions = {
    day: "2-digit",
    month: "short",
    year: "numeric",
  },
) => {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return t("common.notAvailableSymbol");
  }

  return new Intl.DateTimeFormat(locale.value === "ru" ? "ru-RU" : "en-US", options).format(date);
};

export const formatDateTime = (
  value: string | number | Date,
  options: Intl.DateTimeFormatOptions = {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  },
) => formatDate(value, options);

locale.value = resolveInitialLocale();
applyDocumentLocale(locale.value);

export const i18n = {
  locale: readonly(locale),
  getLocale,
  setLocale,
  t,
  formatNumber,
  formatDate,
  formatDateTime,
};
