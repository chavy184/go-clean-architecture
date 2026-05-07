# Go Clean Architecture

> 基于 **整洁架构 (Clean Architecture)** 与 **DDD (领域驱动设计)** 的生产级 Golang 微服务脚手架。

---

## 特性一览

| 特性 | 技术选型 | 说明 |
|------|----------|------|
| 🏗️ 分层架构 | DDD + Clean Architecture | Domain / Application / Delivery / Infrastructure 四层解耦 |
| 💉 依赖注入 | Google Wire | 编译期静态注入，零反射开销 |
| 🌐 多协议交付 | Gin · gRPC · WebSocket | 同一业务逻辑可通过不同协议对外暴露 |
| 🗃️ 数据持久化 | GORM + PostgreSQL | 配置驱动连接，UnitOfWork 事务管理 |
| ⚡ 缓存 | go-redis/v9 | 强类型配置注入 |
| 📝 日志 | Zap + Lumberjack | 依赖注入式日志，自动轮转与归档 |
| 🔍 全链路追踪 | 自研 TraceMiddleware | 自动生成/透传 Trace ID，全层贯穿 Context |
| 📊 可观测性 | Prometheus + pprof | HTTP/gRPC 指标自动采集，在线 CPU/Heap 剖析 |
| 🛡️ 异常捕获 | Recovery 中间件 | HTTP / gRPC / WebSocket 全覆盖 Panic Recovery |
| 🔄 优雅停机 | os/signal | SIGINT/SIGTERM → 排空请求 → 按序关闭资源 |
| 📖 API 文档 | Swagger (swaggo) | 注释即文档，启动即可访问 |
| 🧪 TDD 测试 | 标准库 testing + httptest | 三层 Mock + 红绿重构示例 |
| ⚙️ 配置管理 | Viper | YAML + 环境变量覆盖，多环境隔离 |

---

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose（可选，用于本地中间件）

### 1. 克隆项目

```bash
git clone <your-repo-url> go-clean-architecture
cd go-clean-architecture
```

### 2. 启动依赖中间件

```bash
cd deployments
docker-compose up -d
```

### 3. 配置

```bash
cd backend/config
cp config.default.yaml config.local.yaml
# 编辑 config.local.yaml，填入本地的数据库/Redis 连接信息
```

### 4. 生成 Wire 注入代码

```bash
cd backend
go install github.com/google/wire/cmd/wire@latest
cd cmd/server && wire
```

### 5. 运行

```bash
cd backend
go run cmd/server/main.go
```

启动后可访问：
- **API 服务**：`http://localhost:8080`
- **Swagger 文档**：`http://localhost:8080/swagger/index.html`
- **Prometheus 指标**：`http://localhost:8080/metrics`
- **pprof 性能剖析**：`http://localhost:8080/debug/pprof/`
- **gRPC 服务**：`localhost:9090`

### 6. 运行测试

```bash
cd backend
go test ./internal/... -v -count=1
```

---

## 项目结构

```
project/
├── .docs/              # 架构文档、TDD 指南等
├── api/                # API 契约（protobuf、openapi）
├── deployments/        # Docker / K8s 部署编排
├── backend/            # 后端服务
│   ├── cmd/server/     # 程序入口 + Wire 组装
│   ├── config/         # Viper 配置管理
│   ├── internal/
│   │   ├── domain/     # 领域层（实体、值对象、仓储接口、领域服务）
│   │   ├── application/# 应用层（UseCase、DTO、UoW 接口）
│   │   ├── delivery/   # 交付层（HTTP / gRPC / WebSocket）
│   │   ├── infrastructure/ # 基础设施层（PostgreSQL、Redis、消息队列）
│   │   └── pkg/        # 内部工具库（logger、trace、response）
│   └── test/           # Mock 与测试 Fixture
├── frontend/           # 前端代码
└── README.md
```

> 📚 完整的目录注释与文件说明请参阅 [.docs/architecture.md](.docs/architecture.md)

---

## 核心架构

```
  HTTP / gRPC / WebSocket（Delivery 层）
          │
          ▼
     UseCase（Application 层）
      │            │
      ▼            ▼
  Repository    UnitOfWork     ← 接口定义在 Application 层
  (接口)        (接口)
      │            │
      ▼            ▼
  PostgreSQL     GORM Tx       ← 实现在 Infrastructure 层
```

**依赖规则**：内层不依赖外层。Domain 层是纯净的业务核心，不引入任何框架和数据库包。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [架构设计与实现细节](.docs/architecture.md) | 目录结构规范、9 大核心特性的技术方案详解 |
| [TDD 编写指南](.docs/TDD-Guide.md) | 三层 TDD 实战（红-绿-重构循环、Mock 设计、测试示例） |

---

## 常用命令

```bash
# Wire 依赖注入代码生成
cd backend/cmd/server && wire

# Swagger 文档生成
cd backend && swag init -g cmd/server/main.go --parseDependency --parseInternal

# 运行全部测试（带覆盖率）
cd backend && go test ./... -cover

# 本地开发启动
cd backend && go run cmd/server/main.go
```

---

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go 1.21+ |
| Web 框架 | Gin |
| RPC 框架 | gRPC + protobuf |
| 实时通信 | gorilla/websocket |
| ORM | GORM |
| 数据库 | PostgreSQL |
| 缓存 | Redis (go-redis/v9) |
| 依赖注入 | Google Wire |
| 配置管理 | Viper |
| 日志 | Zap + Lumberjack |
| 监控 | Prometheus + pprof |
| API 文档 | swaggo/swag |
| 测试 | 标准库 testing + httptest |

---

## License

MIT
