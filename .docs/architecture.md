# 架构设计与实现细节

> 本文档详细记录了 go-clean-architecture 的目录结构规范、分层职责以及各核心特性的技术实现方案。

## 目录结构

```
project/
├── .docs/                          # 全栈共享文档（架构决策、API 文档、流程图）
├── api/                            # 对外暴露的 API 契约（前后端/多端共享）
│   ├── protobuf/                   # gRPC 的 .proto 定义文件
│   └── openapi/                    # Swagger/OpenAPI 的 .yaml 规范定义
├── deployments/                    # 基础设施部署与编排
│   ├── docker/                     # 各服务的 Dockerfile
│   │   ├── app.Dockerfile
│   │   └── migrate.Dockerfile
│   ├── k8s/                        # Kubernetes 部署清单
│   └── docker-compose.yml          # 本地开发环境中间件编排（PostgreSQL, Redis, Kafka）
├── backend/                        # 后端服务根目录
│   ├── cmd/
│   │   └── server/
│   │       ├── main.go             # 程序入口，启动HTTP/gRPC/WS，实现优雅停机
│   │       ├── wire.go             # 顶层 Wire 依赖注入组装
│   │       └── wire_gen.go
│   ├── config/
│   │   ├── config.go               # 配置结构体与加载逻辑
│   │   ├── config.default.yaml     # 默认配置模板（脱敏，入库）
│   │   └── config.local.yaml       # 本地配置（git ignore，本地覆盖使用）
│   ├── internal/                   # 私有应用代码（防腐层，核心逻辑）
│   │   ├── domain/                 # 领域层 — 按聚合（Aggregate）划分
│   │   │   ├── user/               # 用户聚合
│   │   │   │   ├── entity.go       # 实体与聚合根
│   │   │   │   ├── valueobject.go  # 值对象
│   │   │   │   ├── repository.go   # 仓储接口定义（不含任何 GORM 或数据库逻辑）
│   │   │   │   └── service.go      # 领域服务
│   │   │   ├── order/              # 订单聚合（示例）
│   │   │   │   ├── entity.go
│   │   │   │   ├── repository.go
│   │   │   │   └── service.go
│   │   │   └── shared/             # 跨聚合共享的领域级概念
│   │   │       ├── event.go        # 领域事件基类
│   │   │       └── error.go        # 领域层自定义错误
│   │   ├── application/            # 应用层 — 用例/应用服务
│   │   │   ├── user/
│   │   │   │   ├── create_user.go  # 业务用例（包含 CreateUserUseCase，注入 uow 与 Repo）
│   │   │   │   ├── get_user.go     # 查询用户用例
│   │   │   │   └── dto.go          # 用例专用输入输出 DTO
│   │   │   └── common/             # 应用层公共组件
│   │   │       ├── uow.go          # 【核心枢纽】UnitOfWork 接口定义（应用层契约）
│   │   │       └── eventbus.go     # 事件总线接口
│   │   ├── delivery/               # 交付层（原 interfaces），负责协议适配
│   │   │   ├── http/               # HTTP (REST/GraphQL) 协议适配
│   │   │   │   ├── v1/             # API 版本控制
│   │   │   │   │   └── user.go     # HTTP Handler，解析请求并调用 Application 层
│   │   │   │   ├── middleware/     # HTTP 中间件集合
│   │   │   │   │   ├── trace.go    # 全链路追踪 (Trace ID 注入)
│   │   │   │   │   ├── recovery.go # 全局异常捕获 (Panic Recovery)
│   │   │   │   │   └── metrics.go  # Prometheus 性能指标打点
│   │   │   │   └── router.go       # Gin 路由注册 (含 pprof 与 /metrics 接口)
│   │   │   ├── grpc/               # gRPC 交付适配器
│   │   │   │   ├── pb/             # 通过 protoc 自动生成的 .pb.go 存放处
│   │   │   │   ├── provider.go     # gRPC 的 Wire ProviderSet
│   │   │   │   ├── server.go       # gRPC Server 的初始化与拦截器配置
│   │   │   │   └── user_service.go # 实现 pb 接口，内部调用 application 层
│   │   │   └── websocket/          # WebSocket 交付适配器
│   │   │       ├── provider.go     # WS 的 Wire ProviderSet
│   │   │       ├── hub.go          # 连接管理器 (管理所有 Client，处理广播/房间)
│   │   │       ├── client.go       # 单个连接的状态机 (包含 ReadPump 和 WritePump)
│   │   │       ├── router.go       # WS 内部的消息路由器 (按消息类型分发)
│   │   │       ├── handler.go      # 接收 WS 消息，调用 application 层
│   │   │       └── pusher_impl.go  # 【关键】实现 application.common.Notifier，反向调用 Hub
│   │   ├── infrastructure/         # 基础设施层 — 技术细节与框架实现
│   │   │   ├── persistence/
│   │   │   │   ├── postgres/       # 数据库底层实现 (GORM)
│   │   │   │   │   ├── user_repo.go# 实现 domain.user.Repository
│   │   │   │   │   ├── uow.go      # 【核心实现】实现 application.common.UnitOfWork (包含 txKey 与 GORM 事务控制)
│   │   │   │   │   └── po.go       # 持久化对象 Persistent Object（带 GORM tag，与 Domain Entity 隔离）
│   │   │   │   └── redis/
│   │   │   ├── messaging/          # 消息队列适配
│   │   │   └── external/           # 外部服务客户端调用
│   │   └── pkg/                    # 仅限当前微服务使用的内部工具库
│   │       ├── logger/
│   │       ├── apperror/
│   │       ├── response/
│   │       └── jwt/
│   ├── pkg/                        # 真正可以跨微服务共享的纯技术包
│   │   └── strutil/                # 无业务逻辑的底层工具
│   ├── test/                       # 测试与集成
│   │   ├── integration/
│   │   ├── mock/                   # Wire 和 UoW 接口的 Mock 文件通常生成在此处
│   │   └── fixture/
│   ├── embedded_scripts/
│   ├── go.mod
│   ├── go.sum
│   └── Makefile                    # 包含 generate, wire, run, test 等指令
├── frontend/                       # 前端代码
└── README.md
```

---

## 核心架构特性与实现细节

本项目严格遵循了**整洁架构 (Clean Architecture)** 与 **领域驱动设计 (DDD)** 的思想，并结合 Golang 生态最佳实践进行了落地实现。

### 1. 依赖注入 (Dependency Injection)
- 全局使用 `github.com/google/wire` 进行依赖注入。
- **按模块注册**：各个层级或功能模块（如基础设施层、领域层、应用层等）均会向外暴露自身的 `ProviderSet` (例如 `postgres.ProviderSet`, `logger.ProviderSet`)。
- **顶层组装**：所有的 `ProviderSet` 最终在 `cmd/server/wire.go` 中进行集中组装，实现了对象创建与业务逻辑的彻底解耦。

### 2. 日志系统 (Logging)
- 采用高性能日志库 `go.uber.org/zap` 替代标准库。
- 采用面向依赖注入的设计，去除了全局日志变量，通过传递 `*zap.Logger` 的形式提供日志能力。
- 结合 `gopkg.in/natefinish/lumberjack.v2` 实现了完善的日志轮转（Log Rotation）机制（支持按文件大小、时间、备份数量自动切割与归档压缩）。

### 3. 数据一致性与事务管理 (Unit of Work)
- 抽象了统一的 `application/common/UnitOfWork` 接口用于控制事务边界，保证了应用层的纯洁性（不包含任何特定数据库的依赖）。
- 基础设施层（PostgreSQL）基于 `gorm.io/gorm` 实现了该接口，并通过上下文 `context.WithValue` 传递私有的 `txKey` 事务句柄。
- 同一请求链路下，底层的各种 Repository 实现均能无缝地从 `context` 中提取同一个数据库事务实例，从而确保跨仓储的数据强一致性。

### 4. Swagger API 接口文档
- 系统使用 `swaggo/swag` 及 `swaggo/gin-swagger` 自动生成接口文档。
- Swagger 注释在 `cmd/server/main.go` 和具体的 Delivery 层（如 `internal/delivery/http/v1/user.go`）中进行标记。
- 可以使用命令 `swag init -g cmd/server/main.go --parseDependency --parseInternal` （需在 backend 目录下执行）生成最新的 `docs` 目录。
- 启动应用后即可访问 `/swagger/index.html` 在线查看并调试 API。

### 5. 基于配置的数据源管理
- 通过 `backend/config/config.go` 集中管理应用级的配置（如 PostgreSQL、Redis 连接等），并且使用强类型的结构体与之映射。
- 采用 **配置驱动实例化** 的模式。PostgreSQL（在 `persistence/postgres/db.go` 中）和 Redis（在 `persistence/redis/cache.go` 中）的连接初始化方法都会要求注入 `*config.Config`，从而获取其配置的 DSN 和连接参数。
- 将连接的初始化函数（如 `NewDB`、`NewRedisClient`）纳入对应的 `ProviderSet` 中，在顶层的 `wire.go` 统一自动编排其依赖顺序和加载流程。

### 6. 全链路追踪与 Context 传递
- **规范 Context 传递**：系统所有核心层的接口方法（包括 Application 层的 UseCase 和 Domain / Infrastructure 层的 Repository）在第一个参数位置强制接收标准库的 `context.Context`，保证上下文信息贯穿整个请求链路。
- **全自动的 Trace ID 注入**：我们在交付层（`delivery/http/middleware/trace.go`）实现了专属中间件。该中间件首先尝试从 HTTP 请求头 (`X-Trace-Id`) 获取前端/调用方传递的 Trace ID。如果未传递，则后端会主动使用 UUID 算法生成。
- **追踪数据透传**：获取或生成的 Trace ID 会通过 `pkg/trace` 工具包，以强类型的私有 context key 的形式被注入到请求的 Context 中。同时我们也会将其自动附加到 HTTP Response Header 返回，极大地方便了跨端排查与微服务场景下的全链路日志串联。

### 7. 优雅停机 (Graceful Shutdown)
- 系统需要监听操作系统的 `SIGINT` 和 `SIGTERM` 信号。
- 接收到信号后，拒绝新的 HTTP/gRPC 请求，等待已在处理的请求完成（或设置一个超时时间如 30 秒）。
- 最后按照依赖顺序安全关闭 PostgreSQL 连接、清理 Redis 连接池、断开所有 WebSocket 客户端，并关闭日志缓冲流（如 `zap.Sync()`），保证正在处理的请求不被粗暴中断。

### 8. 全局异常捕获 (Panic Recovery)
- 在 HTTP 的 Middleware、gRPC 的 Interceptor 以及 WebSocket 的每一个 Goroutine（如 ReadPump/WritePump）中均引入 Recovery 机制。
- 捕获 panic 后，记录包含堆栈信息的 Error 级日志，并向客户端返回统一的 `500 Internal Server Error`，保证 Go 语言的 panic 不会导致整个进程崩溃，同时不影响其他并发请求的正常执行。

### 9. 完整的可观测性 (Observability)
除了基础的 Logging 和 Tracing，系统还集成了 Metrics 和 Profiling：
- **性能监控与指标暴露 (Metrics)**：
  集成 `prometheus/client_golang`，在 HTTP/gRPC 层增加中间件，自动拦截并记录三大核心指标：**QPS**（请求量）、**Latency**（响应延迟直方图）和 **Error Rate**（错误率）。通过 `/metrics` 端点暴露给 Prometheus 抓取，以便后续在 Grafana 中配置大盘和告警。
- **性能剖析 (pprof)**：
  在非公共暴露的内部管理端口（或特定鉴权路由）注册 Go 标准库的 `net/http/pprof` 路由，支持在不停机的情况下实时抓取 CPU Profile 和 Heap Profile，从而为出现 CPU 飙升或内存泄漏时提供在线诊断手段。
