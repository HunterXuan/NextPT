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

Run the full Vue, TypeScript, and server type check before building:

```bash
npm run typecheck
```

Nuxt's build command does not run this check automatically.

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

## Site Branding

Nuxt Node 服务支持在容器启动时配置品牌，无需重新构建镜像：

| 环境变量 | 默认值 | 说明 |
| --- | --- | --- |
| `NUXT_PUBLIC_SITE_NAME` | `NextPT` | 导航、登录注册、页脚和浏览器标题使用的站点名 |
| `NUXT_PUBLIC_SITE_LOGO` | 空 | Logo 路径，例如 `/branding/logo.svg`；未设置或加载失败时使用默认图标 |
| `NUXT_PUBLIC_SITE_LOGO_DARK` | 空 | 暗色背景使用的 Logo 路径；未设置或加载失败时回退到普通 Logo |
| `NUXT_PUBLIC_SITE_FAVICON` | `/favicon.ico` | 浏览器图标路径；挂载自定义文件时建议使用 `/branding/favicon.ico` |

Logo 通过 `<img>` 显示，支持 SVG、PNG、WebP，不会自动修改 SVG 内部颜色。建议提供透明背景、近似正方形的标志，展示区域为 36 × 36 px（页脚 32 × 32 px），图片等比缩放，不会裁切。登录等页面的深色背景区域始终优先使用暗色 Logo。

### Local Testing

仓库保留 `public/branding/.gitkeep`，该目录中的自定义图片已被 Git 和 Docker 构建上下文忽略，不会进入共享镜像。将自己的 `logo.svg`、可选的 `logo-dark.svg` 和 `favicon.ico` 放入 `fe/ui/public/branding/`，然后在 `fe/ui/` 下执行：

```bash
cp .env.example .env
```

编辑 `.env`，取消需要的配置项的注释并填写名称、文件路径：

```dotenv
NUXT_PUBLIC_SITE_NAME=MyPT
NUXT_PUBLIC_SITE_LOGO=/branding/logo.svg
NUXT_PUBLIC_SITE_LOGO_DARK=/branding/logo-dark.svg
NUXT_PUBLIC_SITE_FAVICON=/branding/favicon.ico
```

没有暗色 Logo 时不设置 `NUXT_PUBLIC_SITE_LOGO_DARK`。启动开发服务：

```bash
npm run dev
```

验证步骤：

1. 直接打开 `http://localhost:3000/branding/logo.svg`，确认图片可以访问。
2. 打开首页及 `/login`，确认站点名称、页面标题和 Logo；登录页深色背景区域优先使用暗色 Logo。
3. 在首页切换亮色/暗色主题，检查两个版本的 Logo；未设置暗色版本时复用普通版本。
4. 检查浏览器 favicon，必要时强制刷新或关闭页面后重新打开。
5. 清空 Logo 配置或填写不存在的文件路径，确认回退到默认广播塔图标。

修改 `.env` 后需要重启开发服务；替换图片后刷新页面即可。若使用其它端口，替换以上 URL 中的端口。生产构建预览时，需将测试图片放入 `.output/public/branding/`，并在启动 Node 服务时显式提供环境变量；不要依赖生产服务自动读取 `.env`。

### Kubernetes

`/branding/{file}` 由前端 Node 服务读取挂载文件，支持 SVG、PNG、WebP、ICO，文件名只允许字母、数字、下划线和连字符。这样不会依赖 Nuxt 构建时的静态资源索引。容器工作目录保持 Dockerfile 中的 `/app`；本地开发从 `public/branding/` 读取。以下为 Kustomize 配置示例，图片文件相对于 `kustomization.yaml` 放置：

```yaml
configMapGenerator:
  - name: ui-branding
    files:
      - branding/logo.svg
      - branding/logo-dark.svg
      - branding/favicon.ico
```

在 Deployment 对应容器中添加：

```yaml
env:
  - name: NUXT_PUBLIC_SITE_NAME
    value: MyPT
  - name: NUXT_PUBLIC_SITE_LOGO
    value: /branding/logo.svg
  - name: NUXT_PUBLIC_SITE_LOGO_DARK
    value: /branding/logo-dark.svg
  - name: NUXT_PUBLIC_SITE_FAVICON
    value: /branding/favicon.ico
volumeMounts:
  - name: branding
    mountPath: /app/.output/public/branding/logo.svg
    subPath: logo.svg
    readOnly: true
  - name: branding
    mountPath: /app/.output/public/branding/logo-dark.svg
    subPath: logo-dark.svg
    readOnly: true
  - name: branding
    mountPath: /app/.output/public/branding/favicon.ico
    subPath: favicon.ico
    readOnly: true
```

在 Pod 的 `spec` 下添加：

```yaml
volumes:
  - name: branding
    configMap:
      name: ui-branding
```

没有暗色 Logo 时，省略对应文件、环境变量和挂载项。单个 ConfigMap 数据总大小上限为 1 MiB。保留 Kustomize 默认内容哈希名称，使图片更新触发 Pod 滚动更新；浏览器可能仍缓存同路径图片，需要刷新或更新配置中的资源版本参数（例如 `/branding/logo.svg?v=2`）。挂载文件应在服务启动前准备好。

这些配置只影响前端展示；后端文件中的 `site.name`、`site.url` 仍用于邮件和 TOTP，应在部署时保持站点名称一致。`nuxt generate` 的静态站点不支持容器启动时覆盖配置，需要在生成时设置环境变量。
## Color Themes

Set `NUXT_PUBLIC_SITE_THEME` at deployment time to `teal`, `green`, `sky`,
`indigo`, `violet`, `rose`, `amber`, or `graphite`. The fallback is `sky`.
This is runtime configuration; rebuilding the image is not required.

Users can select a palette in Account Settings > Appearance, independently
of light/dark/system mode. The override is stored in the current browser's
localStorage, not the user account. Choosing the site default removes it.
Unknown theme IDs fall back to `sky`; unavailable browser storage does not
prevent switching themes for the current page.
