// Deprecated: 此文件保留作为 Kratos HTTP 直连方式的参考实现。
// 当前项目已改用 grpc-gateway（参见 grpc_gateway.go），不再使用此文件中的 NewRESTServer。
// 如果将来需要回到 Kratos 原生 HTTP 方式（配合 *_http.pb.go），可参考此文件。
package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/http"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"go-kratos-mcp-demo/cmd/server/assets"
	"go-kratos-mcp-demo/internal/service"

	recommendV1 "go-kratos-mcp-demo/api/gen/go/recommend/service/v1"
)

// NewMiddleware 创建中间件
func newRestMiddleware(
	logger log.Logger,
) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(logger))
	return ms
}

// NewRESTServer new an HTTP server.
// Deprecated: 改用 NewGRPCGatewayServer + grpc-gateway 方案。
func NewRESTServer(
	logger log.Logger,

	recommendService *service.RecommendService,
) *http.Server {

	srv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			newRestMiddleware(logger)...,
		),
	)

	swaggerUI.RegisterSwaggerUIServerWithOption(
		srv,
		swaggerUI.WithTitle("Recommend MCP Server"),
		swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
	)

	// Register routes
	recommendV1.RegisterRecommendServiceHTTPServer(srv, recommendService)

	return srv
}
