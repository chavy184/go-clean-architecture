# TDD 编写指南：Go 整洁架构三层测试实战

> 本文档基于项目 `go-clean-architecture` 的真实代码，演示如何在 **整洁架构 (Clean Architecture)** 的三个核心层中落地 **TDD（测试驱动开发）**。

## 目录

- [TDD 编写指南：Go 整洁架构三层测试实战](#tdd-编写指南go-整洁架构三层测试实战)
  - [目录](#目录)
  - [核心心法](#核心心法)
  - [架构层次与测试策略总览](#架构层次与测试策略总览)
  - [第一层：Infrastructure（仓储实现）](#第一层infrastructure仓储实现)
    - [测试对象](#测试对象)
    - [测试策略](#测试策略)
    - [TDD 循环示例](#tdd-循环示例)
  - [第二层：Application（用例/UseCase）](#第二层application用例usecase)
    - [测试对象](#测试对象-1)
    - [测试策略](#测试策略-1)
    - [TDD 循环示例](#tdd-循环示例-1)
    - [真实案例：TDD 暴露的 Bug](#真实案例tdd-暴露的-bug)
  - [第三层：Delivery（HTTP Handler）](#第三层deliveryhttp-handler)
    - [测试对象](#测试对象-2)
    - [测试策略](#测试策略-2)
    - [TDD 循环示例](#tdd-循环示例-2)
  - [Mock 的设计原则](#mock-的设计原则)
    - [1. Mock 只实现接口，不引入业务逻辑](#1-mock-只实现接口不引入业务逻辑)
    - [2. Mock UoW 直接执行闭包](#2-mock-uow-直接执行闭包)
    - [3. 何时使用 mockgen？](#3-何时使用-mockgen)
  - [运行测试](#运行测试)
    - [测试结果](#测试结果)
  - [最佳实践与避坑指南](#最佳实践与避坑指南)
    - [✅ 推荐做法](#-推荐做法)
    - [❌ 常见错误](#-常见错误)
  - [文件清单](#文件清单)

---

## 核心心法

TDD 的开发循环是一个严格的三步迭代：

```
🔴 红 (Red)     → 先写一个必定失败的测试，明确"我期望什么行为"
🟢 绿 (Green)   → 写最少量的代码让测试通过，不要过度设计
🔵 重构 (Refactor) → 在测试保护下优化代码结构，消除重复
```

> **绝对不要跳过红色阶段。** 如果你写的测试一上来就是绿的，说明你要么测试写得太弱，要么代码已经写好了——这都不是 TDD。

---

## 架构层次与测试策略总览

```
┌─────────────────────────────────────────┐
│  Delivery 层 (HTTP Handler)             │  ← httptest 端到端
│    ↓ 调用                               │
│  Application 层 (UseCase)               │  ← 纯单元测试 + Mock
│    ↓ 调用                               │
│  Infrastructure 层 (Repository 实现)     │  ← 集成测试 / 映射测试
└─────────────────────────────────────────┘
```

| 层级 | 测试类型 | Mock 策略 | 测试文件 |
|------|----------|-----------|----------|
| **Infrastructure** (Repo) | 集成测试 / 单元测试 | 无 Mock，直接测试实现 | `postgres/user_repo_test.go` |
| **Application** (UseCase) | 纯单元测试 | Mock Repo + Mock UoW | `application/user/usecase_test.go` |
| **Delivery** (HTTP Handler) | HTTP 端到端测试 | Mock Repo + Mock UoW + `httptest` | `delivery/http/v1/user_test.go` |

---

## 第一层：Infrastructure（仓储实现）

### 测试对象
`internal/infrastructure/persistence/postgres/user_repo.go` — 实现了 `domain.user.Repository` 接口。

### 测试策略
- 验证 **PO（持久化对象）与 Domain Entity 的映射** 是否正确
- 验证 **构造函数** 返回有效实例
- 对于真正的 SQL 兼容性，建议在 CI 中使用 Docker + PostgreSQL 做集成测试

### TDD 循环示例

**🔴 红：写失败的测试**

```go
// 文件：postgres/user_repo_test.go
func TestUserPO_FieldMapping(t *testing.T) {
    po := postgres.UserPO{
        ID:       "test-001",
        Username: "alice",
    }

    if po.ID != "test-001" {
        t.Errorf("expected ID 'test-001', got '%s'", po.ID)
    }
    if po.Username != "alice" {
        t.Errorf("expected Username 'alice', got '%s'", po.Username)
    }
}
```

> 如果你此时还没定义 `UserPO` 结构体，运行 `go test` 将会编译失败——**这就是红色阶段**。

**🟢 绿：写最少的代码让测试通过**

```go
// 文件：postgres/po.go
type UserPO struct {
    ID       string `gorm:"primaryKey"`
    Username string
}
```

**🔵 重构：** 添加 GORM tag、表名方法等，测试始终保持绿色。

---

## 第二层：Application（用例/UseCase）

### 测试对象
`internal/application/user/create_user.go` 和 `get_user.go`

### 测试策略
- **纯单元测试**：通过注入手写 Mock 完全隔离数据库和框架
- 测试重点：**业务编排逻辑**（入参校验 → 构建实体 → 事务持久化 → 组装响应）

### TDD 循环示例

**🔴 红：空用户名应该被拒绝**

```go
func TestCreateUser_EmptyUsername(t *testing.T) {
    mockRepo := mock.NewMockUserRepository()
    mockUoW := mock.NewMockUnitOfWork()
    uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

    _, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: ""})

    // 期望拿到校验错误
    if err == nil {
        t.Fatal("expected error for empty username, got nil")
    }
    // 断言 Repo 没有被调用（校验应该在持久化之前）
    if mockRepo.SaveCalled {
        t.Error("Save should NOT be called when username is empty")
    }
}
```

> 此时如果 `Execute` 方法体是空的 `return &CreateUserResponse{}, nil`，这个测试**一定会失败**。

**🟢 绿：添加入参校验逻辑**

```go
func (uc *CreateUserUseCase) Execute(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
    if req == nil || req.Username == "" {
        return nil, errors.New("username is required")
    }
    // ...
}
```

**🔴 红：正常创建应该成功**

```go
func TestCreateUser_Success(t *testing.T) {
    mockRepo := mock.NewMockUserRepository()
    mockUoW := mock.NewMockUnitOfWork()
    uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

    resp, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: "alice"})

    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if resp.ID == "" {
        t.Error("expected non-empty ID in response")
    }
    if !mockRepo.SaveCalled {
        t.Error("expected Save to be called")
    }
    if mockRepo.SaveCalledWith.Username != "alice" {
        t.Errorf("expected username 'alice', got '%s'", mockRepo.SaveCalledWith.Username)
    }
}
```

**🟢 绿：补全实体构建和事务逻辑**

```go
func (uc *CreateUserUseCase) Execute(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
    if req == nil || req.Username == "" {
        return nil, errors.New("username is required")
    }

    u := &domainUser.User{
        ID:       uuid.New().String(),
        Username: req.Username,
    }

    err := uc.uow.Do(ctx, func(txCtx context.Context) error {
        return uc.repo.Save(txCtx, u)
    })
    if err != nil {
        return nil, err
    }

    return &CreateUserResponse{ID: u.ID}, nil
}
```

**🔴 红：Repo 出错时应该传播错误**

```go
func TestCreateUser_RepoSaveError(t *testing.T) {
    mockRepo := mock.NewMockUserRepository()
    mockRepo.SaveErr = errors.New("db connection lost")    // 注入故障
    mockUoW := mock.NewMockUnitOfWork()
    uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

    _, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: "alice"})

    if err == nil {
        t.Fatal("expected error when repo fails, got nil")
    }
}
```

**🟢 绿：** 上面的代码已经正确传播了错误——测试直接通过，说明逻辑完备。

> **经验法则**：一个 UseCase 通常需要覆盖以下场景：
> 1. ✅ 入参为空/非法
> 2. ✅ 正常路径（Happy Path）
> 3. ✅ 下游依赖（Repo）出错
> 4. ✅ 事务（UoW）出错
> 5. ✅ 用户不存在（查询类）

### 真实案例：TDD 暴露的 Bug

在开发 `GetUserUseCase` 时，我们先写了一个测试：

```go
func TestGetUserHandler_NotFound(t *testing.T) {
    mockRepo := mock.NewMockUserRepository()     // 空数据
    // ... 构造 Handler 和 router ...

    req := httptest.NewRequest(http.MethodGet, "/api/v1/users/non-existent", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected 500, got %d", w.Code)
    }
}
```

运行后发现 **panic: nil pointer dereference**——因为 `FindByID` 返回了 `nil` 用户，而代码直接访问了 `u.ID`。这就是 TDD 的威力：**测试在你之前发现了 Bug**。

修复方案（🟢 绿）：

```go
u, err := uc.repo.FindByID(ctx, id)
if err != nil {
    return nil, err
}
if u == nil {
    return nil, errors.New("user not found")   // ← TDD 驱动出的防御性代码
}
```

---

## 第三层：Delivery（HTTP Handler）

### 测试对象
`internal/delivery/http/v1/user.go` — HTTP Handler

### 测试策略
- 使用 Go 标准库的 `net/http/httptest` 模拟 HTTP 请求
- 注入 Mock 依赖到 UseCase → Handler 链路
- 验证 **HTTP 状态码** 和 **响应体结构**

### TDD 循环示例

**🔴 红：POST /users 应返回 200**

```go
func TestCreateUserHandler_Success(t *testing.T) {
    mockRepo := mock.NewMockUserRepository()
    mockUoW := mock.NewMockUnitOfWork()
    createUC := appUser.NewCreateUserUseCase(mockUoW, mockRepo)
    getUC := appUser.NewGetUserUseCase(mockRepo)
    handler := v1.NewUserHandler(createUC, getUC)

    body, _ := json.Marshal(appUser.CreateUserRequest{Username: "alice"})
    req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router := gin.New()
    router.POST("/api/v1/users", handler.CreateUser)
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
    }

    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)
    if resp["code"].(float64) != 0 {
        t.Errorf("expected code 0, got %v", resp["code"])
    }
}
```

**🟢 绿：实现 Handler 的请求解析和 UseCase 调用**

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req appUser.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
        return
    }

    resp, err := h.createUC.Execute(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"code": 0, "data": resp})
}
```

**🔴 红：无效 JSON 应返回 400**

```go
func TestCreateUserHandler_BadRequest(t *testing.T) {
    // ... 构造 Handler ...
    req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte("{invalid json")))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }
}
```

---

## Mock 的设计原则

项目中使用**手写 Mock**（位于 `test/mock/mock.go`），遵循以下原则：

### 1. Mock 只实现接口，不引入业务逻辑

```go
type MockUserRepository struct {
    Users          map[string]*user.User  // 内存存储
    SaveCalled     bool                   // 调用断言
    SaveCalledWith *user.User             // 参数捕获
    SaveErr        error                  // 故障注入
}
```

### 2. Mock UoW 直接执行闭包

```go
func (m *MockUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
    if m.DoErr != nil {
        return m.DoErr   // 模拟事务失败
    }
    return fn(ctx)       // 无事务，直接执行
}
```

### 3. 何时使用 mockgen？

| 场景 | 推荐方式 |
|------|----------|
| 接口方法少（≤5个） | 手写 Mock |
| 接口方法多或频繁变更 | `mockgen` 自动生成 |
| 需要精确的调用次数断言 | `gomock` + `mockgen` |

---

## 运行测试

```bash
# 运行全部三层测试（Windows 环境需要指定 GOOS）
$env:GOOS="windows"
go test ./internal/application/user/ \
        ./internal/delivery/http/v1/ \
        ./internal/infrastructure/persistence/postgres/ \
        -v -count=1

# 运行单个层
go test ./internal/application/user/ -v -count=1

# 带覆盖率
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 测试结果

```
=== RUN   TestCreateUser_EmptyUsername        --- PASS
=== RUN   TestCreateUser_NilRequest           --- PASS
=== RUN   TestCreateUser_Success              --- PASS
=== RUN   TestCreateUser_RepoSaveError        --- PASS
=== RUN   TestCreateUser_UoWError             --- PASS
=== RUN   TestGetUser_EmptyID                 --- PASS
=== RUN   TestGetUser_Success                 --- PASS
=== RUN   TestGetUser_RepoError               --- PASS
ok      go-clean-architecture/internal/application/user

=== RUN   TestCreateUserHandler_Success       --- PASS
=== RUN   TestCreateUserHandler_BadRequest    --- PASS
=== RUN   TestGetUserHandler_Success          --- PASS
=== RUN   TestGetUserHandler_NotFound         --- PASS
ok      go-clean-architecture/internal/delivery/http/v1

=== RUN   TestUserPO_FieldMapping             --- PASS
=== RUN   TestNewUserRepository_NotNil        --- PASS
ok      go-clean-architecture/internal/infrastructure/persistence/postgres
```

**14 个测试，全部通过 ✅**

---

## 最佳实践与避坑指南

### ✅ 推荐做法

1. **每个公开方法至少覆盖 3 个测试**：Happy Path + 入参异常 + 依赖故障
2. **测试函数命名使用 `Test<对象>_<场景>` 格式**：如 `TestCreateUser_EmptyUsername`
3. **Mock 中捕获调用参数**：`SaveCalledWith` 让你能断言"传给下游的值是否正确"
4. **在 CI 中用 Docker 跑集成测试**：Repo 层的 SQL 兼容性不能只靠单元测试
5. **表驱动测试（Table-Driven Tests）**：当同一方法有大量边界 case 时，使用 `[]struct` 循环

### ❌ 常见错误

1. **跳过红色阶段**：先写代码再补测试 ≠ TDD，那叫"事后测试"
2. **Mock 中包含业务逻辑**：Mock 越简单越好，复杂 Mock 本身就需要测试
3. **测试依赖执行顺序**：每个 `Test*` 函数必须独立构造自己的 Mock 和状态
4. **忽略错误路径**：只测 Happy Path 是自欺欺人，**故障注入测试（Fault Injection）** 才是 TDD 的灵魂

> **Windows 环境注意**：如果你的 `go env GOOS` 是 `linux`（如 WSL 交叉编译场景），运行测试前需要设置 `$env:GOOS="windows"`，否则测试二进制无法在 Windows 上执行。

---

## 文件清单

| 文件路径 | 作用 |
|----------|------|
| `test/mock/mock.go` | 手写 Mock（MockUserRepository + MockUnitOfWork） |
| `internal/application/user/usecase_test.go` | Application 层 UseCase TDD 测试（8 个用例） |
| `internal/delivery/http/v1/user_test.go` | Delivery 层 HTTP Handler TDD 测试（4 个用例） |
| `internal/infrastructure/persistence/postgres/user_repo_test.go` | Infrastructure 层 Repo TDD 测试（2 个用例） |
