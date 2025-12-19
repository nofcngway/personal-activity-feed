package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nofcngway/feed-service/config"
	server "github.com/nofcngway/feed-service/internal/api/feed_service_api"
	"github.com/nofcngway/feed-service/internal/consumer"
	"github.com/nofcngway/feed-service/internal/pb/feed_api"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(cfg *config.Config, api *server.FeedServiceAPI, cons *consumer.UserActionsConsumer, storage *pgstorage.PGStorage) {
	defer storage.Close()
	defer cons.Close()

	ctx := context.Background()
	go func() {
		backoff := 1 * time.Second
		for {
			err := cons.Consume(ctx)
			if err == nil || ctx.Err() != nil {
				return
			}

			slog.Error("kafka consumer error; retrying", "err", err, "backoff", backoff.String())
			time.Sleep(backoff)
			if backoff < 30*time.Second {
				backoff *= 2
				if backoff > 30*time.Second {
					backoff = 30 * time.Second
				}
			}
		}
	}()

	go func() {
		if err := runGRPCServer(*api, cfg.GRPC.Addr); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %w", err))
		}
	}()

	slog.Info("feed-service starting", "grpc_addr", cfg.GRPC.Addr, "http_addr", cfg.HTTP.Addr)
	if err := runGatewayServer(cfg.HTTP.Addr, cfg.GRPC.Addr, cfg.HTTP.SwaggerPath); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %w", err))
	}
}

func runGRPCServer(api server.FeedServiceAPI, grpcAddr string) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	feed_api.RegisterFeedServiceServer(s, &api)

	slog.Info("gRPC server listening", "addr", grpcAddr)
	return s.Serve(lis)
}

func runGatewayServer(httpAddr, grpcAddr, swaggerPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if swaggerPath != "" {
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			return fmt.Errorf("swagger file not found: %s", swaggerPath)
		}
	}

	r := chi.NewRouter()

	if swaggerPath != "" {
		r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, swaggerPath)
		})
		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		))
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	endpoint := grpcEndpointForGateway(grpcAddr)
	if err := feed_api.RegisterFeedServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening", "addr", httpAddr, "grpc_endpoint", endpoint)
	return http.ListenAndServe(httpAddr, r)
}

func grpcEndpointForGateway(grpcAddr string) string {
	a := strings.TrimSpace(grpcAddr)
	if strings.HasPrefix(a, ":") {
		return "localhost" + a
	}
	return a
}
