# NextPT Frontend

Nuxt frontend for the NextPT private tracker application.

## Scope

The current frontend stage covers:

- homepage
- registration
- login/logout
- user center
- profile update
- password change
- Simplified Chinese, Traditional Chinese, and English UI text

Torrent catalog, forum, and admin screens will be added after the account flow is stable.

## Stack

- Nuxt 4
- Vue 3
- @nuxt/ui
- @nuxtjs/i18n

## Directory

```text
fe/ui/
├── app/             # application source
├── i18n/locales/    # zh-CN, zh-TW, en-US messages
├── public/          # static assets
├── nuxt.config.ts   # Nuxt config
└── package.json
```

## Development

Install dependencies:

```bash
npm install
```

Start local development:

```bash
npm run dev
```

Build for production:

```bash
npm run build
```

Preview production build:

```bash
npm run preview
```

In development, `/api/**` is proxied to `http://localhost:8000/api/**`.
