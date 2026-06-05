export {}

declare module 'vue-i18n' {
  export interface ComposerCustomProperties {}

  export interface Composer {
    locale: import('vue').WritableComputedRef<string>
  }
}
