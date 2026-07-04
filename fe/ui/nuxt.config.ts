// https://nuxt.com/docs/api/configuration/nuxt-config
const nodeEnv = (globalThis as typeof globalThis & {
  process?: {
    env?: {
      NODE_ENV?: string
    }
  }
}).process?.env?.NODE_ENV

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  sourcemap: {
    server: false,
    client: false
  },
  routeRules: {
    // Development-only proxy. Production should route /api through the edge gateway.
    ...(nodeEnv === 'development'
      ? { '/api/**': { proxy: 'http://localhost:8000/api/**' } }
      : {})
  },
  app: {
    head: {
      title: 'NextPT',
      meta: [
        {
          name: 'description',
          content: 'NextPT private tracker community.'
        }
      ]
    },
    pageTransition: { name: 'page', mode: 'out-in' },
    layoutTransition: { name: 'layout', mode: 'out-in' }
  },
  modules: ['@nuxtjs/i18n', '@nuxt/ui'],
  css: ['~/assets/css/main.css'],
  ui: {},
  icon: {
    serverBundle: {
      collections: ['lucide', 'heroicons']
    },
    clientBundle: {
      scan: true
    },
    localApiEndpoint: '/_nuxt_icon'
  },
  i18n: {
    baseUrl: 'http://localhost:3000',
    strategy: 'no_prefix',
    defaultLocale: 'zh-CN',
    langDir: 'locales',
    locales: [
      {
        code: 'zh-CN',
        language: 'zh-CN',
        name: '简体中文',
        file: 'zh-CN.json'
      },
      {
        code: 'zh-TW',
        language: 'zh-TW',
        name: '繁體中文',
        file: 'zh-TW.json'
      },
      {
        code: 'en-US',
        language: 'en-US',
        name: 'English',
        file: 'en-US.json'
      }
    ],
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'nextpt_locale',
      redirectOn: 'root'
    }
  }
})
