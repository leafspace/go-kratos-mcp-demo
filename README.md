# 基于 Go-Kratos 与 MCP 的模块化推荐服务

简短说明：本项目演示如何用 Go-Kratos 框架结合 MCP（模块化协同协议）构建可扩展、可观测的推荐服务，包含召回 / 过滤 / 排序等模块化流程与测试范例。

## 核心特性

- 标准化契约：以 Protobuf 定义服务与上下文，便于多语言客户端与版本管理。
- 模块化协同：通过 MCP 将召回、排序、过滤等模块串联，模块可独立开发与部署。
- 多协议支持：同一份 proto 同时暴露 gRPC、grpc-gateway HTTP 与 MCP Tool，满足不同调用场景。
- 配置驱动：服务地址、MCP 类型等通过 YAML 配置注入，无需硬编码。
- 可观测性与测试：支持链路追踪、日志、指标，并提供端到端与单元测试示例（兼容 Protobuf JSON 中 int64 被编码为 string 的情况）。

## 架构概览

- 请求流：客户端 → gRPC / HTTP / MCP 请求 → 统一 service 方法 → 召回 → 排序 → 过滤 → 返回 Protobuf 响应或 MCP Tool 结果。
- 依赖注入：使用 Kratos 的 wire 管理模块依赖，s/recall/rank/filter 服务独立实现逻辑并通过接口解耦。
- 协议入口：Gateway 地址（如 `:8000`）使用 h2c 按 `Content-Type: application/grpc` 分流 gRPC 与 grpc-gateway HTTP；MCP 地址（如 `:8080`）提供 MCP Streamable HTTP。
- 工具注册：MCP 服务启动时调用 `protoc-gen-go-mcp` 生成的 `Register*ServiceMCPHandler`，Tool Schema 与 handler 由 proto 生成。

## 配置说明

配置文件 `configs/recommend.yaml`：

```yaml
name: recommend-mcp-server
version: v1.0.0

gateway:
  addr: ":8000"        # gRPC + HTTP Gateway 地址

mcp:
  addr: ":8080"        # MCP 服务地址
  type: "http"         # MCP 服务类型: stdio / sse / http
  name: "Recommend MCP Server"
  version: "1.0.0"
  server:
    start_timeout: "15s"
  context:
    max_history_len: 10

service:
  recall:
    top_k: 50
  rank:
    top_k: 10
  filter:
    blacklist:
      - 999
      - 888
```

## 快速上手

1. 安装生成器：`go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest`，并安装本地 `github.com/leafspace/protoc-gen-go-mcp` fork。
2. 生成配置代码：`cd backend && make config`
3. 生成 proto 代码：`make api`
4. 生成依赖注入：`make wire`
5. 启动服务：`make run`（同时启动 MCP 与 gRPC/gateway 服务）
6. 本地测试：
   - HTTP 调用：`POST http://127.0.0.1:8000/api/v1/recommend`
   - MCP 调用：`POST http://127.0.0.1:8080/mcp`

## Proto 统一契约

本 demo 的核心思路是让 `proto` 成为 gRPC、HTTP 与 MCP Tool 的统一服务契约：

- 统一 proto 契约：业务只维护一份 `*.proto`，同时生成 `*.pb.go`、`*_grpc.pb.go`、`*_http.pb.go`、`*.pb.gw.go` 与 `*.pb.mcp.go`，减少 HTTP、gRPC、MCP 三套接口重复定义。
- gRPC 与 HTTP 兼容：通过 `grpc-gateway` 生成网关代码，同一个 service 实现可以同时服务 gRPC 调用与 REST/HTTP 调用。
- MCP Tool 自动生成：通过 `protoc-gen-go-mcp` 从 proto service/rpc/message 自动生成 Tool、参数 Schema、请求解码、响应编码与 `Register*ServiceMCPHandler` 注册函数。
- Tool 注册归一：MCP 服务启动时只需要调用生成的 `Register*ServiceMCPHandler`，工具自动挂载。
- 字段描述进入 Schema：使用 `(mcp.v1.desc)` 注解，MCP Tool 参数 Schema 同时包含字段类型与描述。

推荐的 proto 写法：

```proto
import "google/api/annotations.proto";
import "mcp/v1/annotations.proto";

service RecommendService {
  rpc TriggerRecommend(RecommendRequest) returns (RecommendResponse) {
    option (google.api.http) = {
      post: "/api/v1/recommend"
      body: "*"
    };
  }
}

message RecommendRequest {
  UserActionContext action_ctx = 1 [(mcp.v1.desc) = "用户行为上下文"];
}
```

## 测试与注意点

- 默认单元测试不依赖已启动的 MCP 服务；`TestMcpClient` 是集成测试，需要先启动服务后手动打开。
- 推荐在测试中对 responseTime 使用 time.RFC3339Nano 校验；对 historyItemIds、blacklistItemIds 等 ID 列表做兼容字符串/数字的断言。
- 在断言数值时，建议编写通用转换函数以支持 float64/string/int64 三种表现形式。

## 后续可优化方向

1. 将 `protoc-gen-go-mcp` 的 `MCP` 后缀能力上游化到 fork，移除 Makefile 中的 perl 后处理命令。
2. 将 Swagger UI 挂到 grpc-gateway mux 上，保留原 demo 的在线接口文档能力。
3. 为 MCP 注册增加端到端测试，启动 in-memory 或随机端口服务后校验工具列表与 schema。
4. gRPC 服务增加 TLS / 认证中间件（当前使用 insecure）。
5. 将 MCP context 的 `ExtraProperties` 能力用于链路追踪信息透传。

## 结论

本项目示例展示了「框架 + 协议」的协同价值：Kratos 提供服务治理与依赖注入能力，MCP 提供模块化协同标准。对需要模块化、可扩展且具备跨服务契约的推荐系统，该方案适合生产化落地。

## 项目链接

- GitHub项目地址: <https://github.com/tx7do/go-kratos-mcp-demo>
- Gitee项目地址: <https://gitee.com/tx7do/go-kratos-mcp-demo>
- MCP 协议封装库: <https://github.com/mark3labs/mcp-go>
- Kratos Transport MCP 扩展: [github.com/tx7do/kratos-transport/transport/mcp](https://github.com/tx7do/kratos-transport/tree/main/transport/mcp)
- protoc-gen-go-mcp (leafspace fork): <https://github.com/leafspace/protoc-gen-go-mcp>
