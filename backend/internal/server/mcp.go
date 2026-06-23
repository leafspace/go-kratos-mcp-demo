package server

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"

	mcpRuntime "github.com/leafspace/protoc-gen-go-mcp/pkg/runtime"
	mark3Mcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	mcpServer "github.com/tx7do/kratos-transport/transport/mcp"

	"go-kratos-mcp-demo/api/gen/go/conf"
	recommendV1 "go-kratos-mcp-demo/api/gen/go/recommend/service/v1"
	"go-kratos-mcp-demo/internal/service"
)

func NewMcpServer(cfg *conf.RecommendConfig, _ log.Logger, recommendService *service.RecommendService) *mcpServer.Server {
	mcpCfg := cfg.GetMcp()

	addr := mcpCfg.GetAddr()
	if addr == "" {
		addr = ":8080"
	}
	srvType := mcpServer.ServerType(mcpCfg.GetType())
	if srvType == "" {
		srvType = mcpServer.ServerTypeHTTP
	}
	srvName := mcpCfg.GetName()
	if srvName == "" {
		srvName = "Recommend MCP Server"
	}
	srvVersion := mcpCfg.GetVersion()
	if srvVersion == "" {
		srvVersion = "1.0.0"
	}

	srv := mcpServer.NewServer(
		mcpServer.WithServerName(srvName),
		mcpServer.WithServerVersion(srvVersion),
		mcpServer.WithMCPServeType(srvType),
		mcpServer.WithMCPServeAddress(addr),
		mcpServer.WithMCPServerOptions(
			server.WithToolCapabilities(false),
			server.WithRecovery(),
		),
	)

	recommendV1.RegisterRecommendServiceMCPHandler(wrapMcpServer(srv), recommendService)

	return srv
}

type runtimeMcpServer struct {
	srv *mcpServer.Server
}

func wrapMcpServer(srv *mcpServer.Server) mcpRuntime.MCPServer {
	return &runtimeMcpServer{srv: srv}
}

// AddTool is called by protoc-gen-go-mcp generated Register*Handler functions.
func (s *runtimeMcpServer) AddTool(tool mcpRuntime.Tool, handler mcpRuntime.ToolHandler) {
	mcpTool := mark3Mcp.Tool{
		Name:            tool.Name,
		Description:     tool.Description,
		RawInputSchema:  json.RawMessage(tool.RawInputSchema),
		RawOutputSchema: json.RawMessage(tool.RawOutputSchema),
	}

	if err := s.srv.RegisterHandler(mcpTool, func(ctx context.Context, request mark3Mcp.CallToolRequest) (*mark3Mcp.CallToolResult, error) {
		result, err := handler(ctx, &mcpRuntime.CallToolRequest{
			Arguments: request.GetArguments(),
		})
		if err != nil {
			return nil, err
		}
		if result == nil {
			return nil, nil
		}
		if result.IsError {
			return mark3Mcp.NewToolResultError(result.Text), nil
		}

		mcpResult := mark3Mcp.NewToolResultText(result.Text)
		mcpResult.StructuredContent = result.StructuredContent
		return mcpResult, nil
	}); err != nil {
		log.Errorf("failed to register generated MCP tool %s: %v", tool.Name, err)
	}
}
