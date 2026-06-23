package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"

	"go-kratos-mcp-demo/api/gen/go/conf"
	recommendV1 "go-kratos-mcp-demo/api/gen/go/recommend/service/v1"
)

func TestRecommendService_TriggerRecommend(t *testing.T) {
	// 准备依赖服务
	logger := log.DefaultLogger
	rc := testRecommendConfig()
	recallService := NewRecallService(logger, rc)
	rankService := NewRankService(logger, rc)
	filterService := NewFilterService(logger, rc)

	// 创建推荐服务
	service := NewRecommendService(logger, recallService, rankService, filterService)

	request := &recommendV1.RecommendRequest{
		ActionCtx: &recommendV1.UserActionContext{
			UserFeature: &recommendV1.UserFeature{
				UserId: "user123",
			},
			TriggerItemId: 101,
			Scene:         recommendV1.SceneType_SCENE_TYPE_HOME,
			ActionType:    recommendV1.ActionType_ACTION_TYPE_CLICK,
		},
	}

	// 执行测试
	ctx := context.Background()
	result, err := service.TriggerRecommend(ctx, request)

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.GetRequestId())
	assert.Equal(t, recommendV1.ContextStage_CONTEXT_STAGE_OUTPUT, result.GetOutput().GetStage())
	assert.NotEmpty(t, result.GetOutput().GetRecommendItems())
}

func TestMcpClient(t *testing.T) {
	t.Skip("integration test: start the demo server before running")

	ctx := context.Background()

	httpTransport, err := transport.NewStreamableHTTP("http://localhost:8080/mcp")
	assert.NoError(t, err)
	assert.NotNil(t, httpTransport)

	mcpClient := client.NewClient(
		httpTransport,
	)
	assert.NotNil(t, mcpClient)
	defer mcpClient.Close()

	err = mcpClient.Start(ctx)
	assert.NoError(t, err)

	initRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    "recommend-http-client",
				Version: "1.0.0",
			},
		},
	}

	_, err = mcpClient.Initialize(ctx, initRequest)
	assert.NoError(t, err)

	// 调用推荐工具
	result, err := mcpClient.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "recommend",
			Arguments: map[string]interface{}{
				"userFeature": map[string]interface{}{
					"userId": "user123",
				},
				"triggerItemId": int64(1001),
				"scene":         "homepage",
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Content)

	if len(result.Content) == 0 {
		t.Fatalf("empty result.Content")
	}
	for i, c := range result.Content {
		switch ct := c.(type) {
		case mcp.TextContent:
			var data map[string]interface{}
			if err = json.Unmarshal([]byte(ct.Text), &data); err != nil {
				t.Fatalf("unmarshal text content[%d] failed: %v", i, err)
			}
			assert.NotEmpty(t, data)
			t.Logf("text content[%d]=%#v", i, data)

			root := data
			assert.Equal(t, "CONTEXT_STAGE_OUTPUT", fmt.Sprintf("%v", root["stage"]))
			assert.NotEmpty(t, root["requestId"])

			inputCtx, ok := root["inputCtx"].(map[string]interface{})
			if !ok {
				t.Fatalf("missing inputCtx")
			}

			if bl, ok := inputCtx["blacklistItemIds"].([]interface{}); ok {
				assert.True(t, len(bl) >= 0)
				assert.Equal(t, bl[0], "999")
				assert.Equal(t, bl[1], "888")
			} else {
				t.Fatalf("blacklistItemIds missing or invalid")
			}

		default:
			t.Logf("content[%d] unexpected type=%T value=%#v", i, c, c)
		}
	}

	t.Logf("完整结果: %#v", result)
}

func testRecommendConfig() *conf.RecommendConfig {
	return &conf.RecommendConfig{
		Service: &conf.RecommendConfig_Service{
			Recall: &conf.RecommendConfig_Recall{TopK: 50},
			Rank:   &conf.RecommendConfig_Rank{TopK: 10},
			Filter: &conf.RecommendConfig_Filter{Blacklist: []int64{999, 888}},
		},
	}
}
