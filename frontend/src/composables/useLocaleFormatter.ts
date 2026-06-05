import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const LOCALE_TAGS: Record<string, string> = {
  ru: 'ru-RU',
  kz: 'kk-KZ',
  en: 'en-US',
}

export function useLocaleFormatter() {
  const { locale } = useI18n()
  const dateLocaleTag = computed(() => LOCALE_TAGS[locale.value] ?? 'ru-RU')

  function fmtDate(iso: string, options?: Intl.DateTimeFormatOptions) {
    return new Date(iso).toLocaleDateString(dateLocaleTag.value, options)
  }
  function fmtTime(iso: string, options?: Intl.DateTimeFormatOptions) {
    return new Date(iso).toLocaleTimeString(
      dateLocaleTag.value,
      options ?? { hour: '2-digit', minute: '2-digit' },
    )
  }
  function fmtDateTime(iso: string, options?: Intl.DateTimeFormatOptions) {
    return new Date(iso).toLocaleString(dateLocaleTag.value, options)
  }

  return { locale, dateLocaleTag, fmtDate, fmtTime, fmtDateTime }
}
