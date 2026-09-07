# Server Resources

服务端静态资源通过 `/api/public/**` 提供。当前目录仅保留结构占位，实际站点品牌和分类图片由部署或项目资源补充。

```text
resource/
├── public/
│   └── images/
│       ├── brand/        # Logo、favicon 等站点品牌资源
│       └── categories/   # 种子分类图标
└── template/             # 服务端模板资源预留目录
```

不要将私密配置、用户上传内容或对象存储凭据放入此目录。
