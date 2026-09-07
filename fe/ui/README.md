# NextPT Frontend

NextPT 的 Nuxt 4 前端，覆盖普通用户站点与管理后台。当前已接入账号与安全设置、种子和论坛、求种与续种、站内通知、站务信箱、在线聊天、用户任务、魔力商城和后台管理页面。

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

前置条件：Node.js 24。安装依赖：

```bash
npm ci
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

开发模式下，`/api/**` 会代理到 `http://localhost:8000/api/**`。生产构建不内置 API 代理，需由网关或反向代理将该路径转发至后端服务。

容器构建使用 [Dockerfile](Dockerfile) 生成 Nuxt Node 服务镜像；镜像标签和多架构发布由仓库根目录的 GitHub Actions 工作流负责。
