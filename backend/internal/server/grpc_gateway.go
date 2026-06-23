package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"go-kratos-mcp-demo/api/gen/go/conf"
	recommendV1 "go-kratos-mcp-demo/api/gen/go/recommend/service/v1"
	"go-kratos-mcp-demo/internal/service"
)

type GRPCGatewayServer struct {
	addr    string
	httpMux *runtime.ServeMux
	grpc    *kgrpc.Server
}

func NewGRPCGatewayServer(cfg *conf.RecommendConfig, logger log.Logger, recommendService *service.RecommendService) *GRPCGatewayServer {
	addr := cfg.GetGateway().GetAddr()
	if addr == "" {
		addr = ":8000"
	}

	grpcServer := kgrpc.NewServer()
	recommendV1.RegisterRecommendServiceServer(grpcServer, recommendService)

	gatewayMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := recommendV1.RegisterRecommendServiceHandlerFromEndpoint(context.Background(), gatewayMux, addr, opts); err != nil {
		log.NewHelper(logger).Fatalf("failed to register grpc-gateway handler: %v", err)
	}

	return &GRPCGatewayServer{
		addr:    addr,
		httpMux: gatewayMux,
		grpc:    grpcServer,
	}
}

func (s *GRPCGatewayServer) Start(context.Context) error {
	return http.ListenAndServe(s.addr, grpcGatewayHandler(s.grpc, s.httpMux))
}

func (s *GRPCGatewayServer) Stop(context.Context) error {
	return nil
}

func grpcGatewayHandler(grpcServer *kgrpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
			return
		}
		otherHandler.ServeHTTP(w, r)
	}), &http2.Server{})
}
