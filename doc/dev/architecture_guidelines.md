# Logic 架构设计规范 (CQRS / Usecase-Domain 分离)

NextPT 项目后端为了解决 Web 请求上下文（`context`）污染底层业务逻辑，以及解决业务入口逻辑与纯表 CRUD 操作耦合过深的问题，确立了 **Usecase / Domain** 的两层分级，以及基于 `model.Actor` 的全局身份显式传递机制。后续所有模块（包括老模块重构与新功能开发）均应参考本规范。

## 一、核心设计概念

### 1. `Actor` 全局身份模型
原先从上下文中隐式获取用户信息的 `contexts.GetUser(ctx)` 已被彻底废弃。
所有经过 IAM (Identity and Access Management) 模块认证的请求，其中间件会解析出用户信息并存入 `model.Actor` 结构体中。

- **获取限制**：`Controller`（接入层）是**唯一**允许使用 `contexts.GetActor(ctx)` 从上下文获取身份的地方。
- **显式传递**：`Controller` 获取到 `Actor` 后，必须将其作为**独立参数**显式地向下级 `Usecase` 传递，剥离 Usecase 对特定 Web 框架上下文读取的依赖。

### 2. Usecase 业务用例层 (`xxx_usecase.go`)
本层是任何具体业务动作的“指挥官”或“编排者”。

- **核心职责**：
  - 接收 Controller 传来的业务参数和 `Actor`。
  - 执行综合业务策略判断（前置条件校验）。
  - 编排下层一个或多个 Domain 的基础能力，推进业务流转。
  - 负责跨 Domain、跨资源的事务编排；这类事务应由 Usecase 开启。
- **禁忌**：
  - 绝对不可以在此层包含纯粹的数据库表插入、更新语句，或执行复杂的 QueryBuilder 拼接。
  - 不得使用 `contexts.GetActor(ctx)`。

### 3. Domain 领域基础层 (`xxx_domain.go`)
本层是底层基础能力的提供者，是具体表操作和原子校验动作的集合。

- **核心职责**：
  - 提供单表的增、删、改、查等原子方法（例如 `GetTopicById`, `InsertTopic`）。
  - 提供不依赖大环境上下文的权限校验等静态策略检查方法（例如 `CheckTopicWritePolicy`）。
- **禁忌**：
  - 通常直接使用 `Ctx(ctx)` 执行语句，若外层 Usecase 开启了事务则自动加入。仅当单个 Domain 方法需要保证自身单资源读改写原子性时，可以在该方法内部开启紧凑事务；不得在其中编排跨 Domain 的业务流程。
  - 绝对不允许在此层内自发地跨界调用其它 Domain 去组合宏大的业务流（跨 Domain 调用的编排职责归 Usecase 所有）。
  - 不得使用 `contexts.GetActor(ctx)`。

---

## 二、代码范例 (以 Forum 模块为例)

### Controller (入口：负责接收与派发)
```go
func (c *ControllerV1) TopicCreate(ctx context.Context, req *v1.TopicCreateReq) (res *v1.TopicCreateRes, err error) {
	// 获取 Actor 并显式传给 Usecase
	actor := contexts.GetActor(ctx)
	id, err := service.ForumTopicUsecase().Create(ctx, actor, req.TopicCreateInp)
	if err != nil {
		return nil, err
	}
	return &v1.TopicCreateRes{TopicCreateOut: forumout.TopicCreateOut{Id: id}}, nil
}
```

### Usecase (编排：负责校验、事务与组织流程)
```go
func (s *sForumTopicUsecase) Create(ctx context.Context, actor *model.Actor, in forumin.TopicCreateInp) (uint64, error) {
	// 1. 基础鉴权
	if actor == nil {
		return 0, gerror.New(gi18n.T(ctx, "forum.general.unauthorized"))
	}
	
	// 2. 检查前置依赖是否存在
	var node entity.ForumNode
	if err := dao.ForumNode.Ctx(ctx).Where(dao.ForumNode.Columns().Id, in.NodeId).Scan(&node); err != nil || node.Id == 0 {
		return 0, gerror.New(gi18n.T(ctx, "forum.node.not_found"))
	}

	// 3. 权限策略检查 (借用 Domain 层原子能力)
	if err := service.ForumNodeDomain().CheckNodeWritePolicy(ctx, actor, &node); err != nil {
		return 0, err
	}

	// 4. 开启事务，编排下游写入操作
	var topicId uint64
	err := dao.ForumTopic.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := service.ForumTopicDomain().InsertTopic(ctx, actor, in)
		if err != nil {
			return err
		}
		topicId = id

		// 这里可继续调用 NodeDomain 或 UserDomain 更新其它统计数据...
		return nil
	})
    
	if err != nil {
		return 0, gerror.Wrap(err, "failed to create topic")
	}
	return topicId, nil
}
```

### Domain (干活：纯粹的数据读写及静态逻辑判定)
```go
// 策略校验
func (s *sForumNodeDomain) CheckNodeWritePolicy(ctx context.Context, actor *model.Actor, node *entity.ForumNode) error {
	userLevel := 0
	if actor != nil {
		userLevel = actor.RoleLevel
	}
	if userLevel < int(node.MinRoleWrite) {
		return gerror.New(gi18n.T(ctx, "forum.node.write_permission_denied"))
	}
	return nil
}

// 纯粹写入
func (s *sForumTopicDomain) InsertTopic(ctx context.Context, actor *model.Actor, in forumin.TopicCreateInp) (uint64, error) {
	id, err := dao.ForumTopic.Ctx(ctx).Data(entity.ForumTopic{
		NodeId:      in.NodeId,
		UserId:      actor.Id,
		Subject:     in.Subject,
		Content:     in.Content,
	}).InsertAndGetId()
	return uint64(id), err
}
```

---

## 三、国际化 (i18n) 最佳实践
- 所有提供给客户端的错误信息**绝不能**直接硬编码中英文字符串。
- **所有**能够抛出用户可见错误的方法签名中必须包含 `ctx context.Context`。
- 使用 `gerror.New(gi18n.T(ctx, "i18n.key"))` 包裹错误。
- 去 `manifest/i18n/` 目录下（如 `en-US`, `zh-CN`, `zh-TW`），补齐对应的 `yaml` 文件与对应的 Key 值翻译。
- **YAML 文件维护规范**：
  1. 必须使用带点的扁平键格式（如 `admin.config.group_req: 配置组不能为空`），绝不允许使用嵌套字典结构。
  2. 翻译词条必须按 Key 字母顺序排序。
  3. 必须根据 `domain.resource` 的前缀格式进行逻辑分块（如 `admin.config` 为一块，`admin.forum` 为另一块）。
  4. 不同的模块块之间必须加一个空行，以保证最佳的阅读和维护体验。

---

## 四、Import 代码规范
为了保持代码整洁和一致性，所有的非 `gf` 工具自动生成的可编辑 `.go` 文件的 `import` 块必须严格分为三段，并且严格按照以下顺序排列，每个分块之间留一个空行：
1. **Golang 原生依赖** (e.g. `context`, `fmt`, `strings`)
2. **服务本身代码依赖** (e.g. `server/api/...`, `server/internal/...`)
3. **外部第三方依赖** (e.g. `github.com/gogf/gf/v2/...`)

## 五、参数传递风格规范
在 `logic` 下辖的 `usecase`（业务编排层）与 `domain`（领域基础层）的方法签名中，必须严格遵循以下参数传值规范：
1. **输入参数 (Inp) 必须使用值传递**：例如 `in adminin.UserListInp` 而不是 `in *adminin.UserListInp`。这能避免由于传引用而在深层逻辑中被意外篡改（保证数据不变性），且小结构体在栈上分配成本极低，有助于减轻 GC 压力。
2. **输出参数 (Out) 必须使用指针返回**：例如 `(*adminout.UserListOut, error)`。因为业务逻辑常常需要同时返回结果和 `error`，在发生错误时返回 `nil, err` 会比返回空壳结构体语义更清晰，同时避免了大结构体在返回时的内存拷贝开销。
3. **Domain 层的方法参数设计**：
   - 尽量直接接收基础数据类型（如 `userId uint64`, `subject string`）或数据库实体（如 `topic *entity.ForumTopic`），让 Domain 保持高度原子性和可复用性。
   - 当参数较多时，可以直接复用并接收 `usecase` 传下来的值类型 `Inp` 结构体（如 `in forumin.TopicCreateInp`）。
   - 必须显式接收 `actor *model.Actor`（如果逻辑需要身份信息），绝对禁止在 Domain 内部读取上下文隐式凭证（`contexts.GetActor`）。

---

## 六、API 路由声明规范
在定义 API 的 `g.Meta` 时，必须严格遵循 OpenAPI 规范：
1. **路径参数（Path Parameter）必须使用 `{id}` 格式**：例如 `path:"/users/{id}"`，绝对禁止使用 `:id` 的传统框架格式（如 `/users/:id`），以确保自动生成的 Swagger/OpenAPI 文档的规范性和一致性。
2. **自定义动作（Custom Action）允许使用 `:` 后缀**：例如 `path:"/users/{id}:ban"`，这是符合 RESTful 及 OpenAPI 规范的标准动词扩展用法。
3. **结构体初始化**：在 Controller 中进行返回结构体初始化时，禁止使用 Unkeyed 格式（如 `&v1.xxxRes{*out}`），必须使用 Keyed 具名赋值（如 `&v1.xxxRes{XxxOut: *out}`），以保持代码规范，消除 Go Linter 警告。
