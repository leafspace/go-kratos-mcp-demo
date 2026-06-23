package service

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	recommendV1 "go-kratos-mcp-demo/api/gen/go/recommend/service/v1"
)

type RecommendService struct {
	recommendV1.RecommendServiceHTTPServer
	recommendV1.UnimplementedRecommendServiceServer

	recallService *RecallService
	rankService   *RankService
	filterService *FilterService

	userHistoryCache map[string][]int64
	maxHistoryLen    int

	log *log.Helper

	sync.RWMutex
}

func NewRecommendService(
	logger log.Logger,
	recallService *RecallService,
	rankService *RankService,
	filterService *FilterService,
) *RecommendService {
	l := log.NewHelper(log.With(logger, "module", "file/service/mcp-service"))

	return &RecommendService{
		recallService: recallService,
		rankService:   rankService,
		filterService: filterService,

		userHistoryCache: make(map[string][]int64),
		maxHistoryLen:    10,

		log: l,
	}
}

// TriggerRecommend 推荐接口
func (s *RecommendService) TriggerRecommend(ctx context.Context, req *recommendV1.RecommendRequest) (*recommendV1.RecommendResponse, error) {
	actionCtx := req.GetActionCtx()
	finalOutput, requestId, err := s.runRecommend(ctx, actionCtx)
	if err != nil {
		return nil, err
	}

	return &recommendV1.RecommendResponse{
		Output:    finalOutput,
		RequestId: requestId,
	}, nil
}

func (s *RecommendService) runRecommend(ctx context.Context, actionCtx *recommendV1.UserActionContext) (*recommendV1.RecommendOutput, string, error) {
	if actionCtx == nil {
		actionCtx = &recommendV1.UserActionContext{}
	}

	requestId := uuid.New().String()
	actionCtx.RequestId = requestId
	actionCtx.Stage = recommendV1.ContextStage_CONTEXT_STAGE_INPUT
	if actionCtx.GetActionTime() == nil {
		actionCtx.ActionTime = timestamppb.Now()
	}

	// 更新用户历史
	userId := actionCtx.GetUserFeature().GetUserId()
	triggerItemId := actionCtx.GetTriggerItemId()

	s.Lock()
	history := s.userHistoryCache[userId]
	history = append([]int64{triggerItemId}, history...)
	if len(history) > s.maxHistoryLen {
		history = history[:s.maxHistoryLen]
	}
	s.userHistoryCache[userId] = history
	s.Unlock()

	// 召回
	recallInput := &recommendV1.RecallInputContext{
		Stage:          recommendV1.ContextStage_CONTEXT_STAGE_RECALL,
		RequestId:      requestId,
		UserFeature:    actionCtx.GetUserFeature(),
		Scene:          actionCtx.GetScene(),
		HistoryItemIds: history,
		RecallTopK:     s.recallService.rc.GetService().GetRecall().GetTopK(),
	}
	recallOutput, err := s.recallService.Recall(ctx, recallInput)
	if err != nil {
		s.log.Errorw(ctx, "recall failed", "err", err)
		return nil, "", recommendV1.ErrorInternalServerError("召回失败")
	}

	// 排序
	rankInput := &recommendV1.RankInputContext{
		Stage:     recommendV1.ContextStage_CONTEXT_STAGE_RANK,
		RequestId: requestId,
		RecallCtx: recallOutput,
		RankTopK:  s.rankService.rc.GetService().GetRank().GetTopK(),
	}
	rankOutput, err := s.rankService.Rank(ctx, rankInput)
	if err != nil {
		s.log.Errorw(ctx, "rank failed", "err", err)
		return nil, "", recommendV1.ErrorInternalServerError("排序失败")
	}

	// 过滤
	filterInput := &recommendV1.FilterInputContext{
		Stage:            recommendV1.ContextStage_CONTEXT_STAGE_FILTER,
		RequestId:        requestId,
		RankCtx:          rankOutput,
		BlacklistItemIds: s.filterService.rc.GetService().GetFilter().GetBlacklist(),
		PurchasedItemIds: []int64{},
	}
	finalOutput, err := s.filterService.Filter(ctx, filterInput)
	if err != nil {
		s.log.Errorw(ctx, "filter failed", "err", err)
		return nil, "", recommendV1.ErrorInternalServerError("过滤失败")
	}

	return finalOutput, requestId, nil
}

// GetRecommendHistory 获取推荐历史(待实现)
func (s *RecommendService) GetRecommendHistory(ctx context.Context, req *recommendV1.RecommendHistoryRequest) (*recommendV1.RecommendHistoryResponse, error) {
	return &recommendV1.RecommendHistoryResponse{
		Total:    0,
		Page:     req.GetPage(),
		PageSize: req.GetPageSize(),
	}, nil
}
