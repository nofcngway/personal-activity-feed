package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nofcngway/auth-action-service/config"
	server "github.com/nofcngway/auth-action-service/internal/api/auth_action_service_api"
	kafkaproducer "github.com/nofcngway/auth-action-service/internal/kafka/producer"
	auth_pb "github.com/nofcngway/auth-action-service/internal/pb/auth_action_api"
	"github.com/nofcngway/auth-action-service/internal/sessions"
	"github.com/nofcngway/auth-action-service/internal/storage/pgstorage"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(cfg *config.Config, api *server.AuthActionServiceAPI, storage *pgstorage.PGStorage, sess *sessions.RedisStore, producer *kafkaproducer.Producer) {
	defer storage.Close()
	defer producer.Close()
	defer sess.Close()

	go func() {
		if err := runGRPCServer(*api, cfg.GRPC.Addr); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %w", err))
		}
	}()

	slog.Info("auth-action-service starting", "grpc_addr", cfg.GRPC.Addr, "http_addr", cfg.HTTP.Addr)
	if err := runGatewayServer(cfg.HTTP.Addr, cfg.GRPC.Addr, cfg.HTTP.SwaggerPath); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %w", err))
	}
}

func runGRPCServer(api server.AuthActionServiceAPI, grpcAddr string) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	auth_pb.RegisterAuthActionServiceServer(s, &api)

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
			serveSwaggerJSONWithAuth(w, swaggerPath)
		})

		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		))
	}

	// Прокидываем заголовки в metadata, чтобы gRPC-методы могли читать Authorization
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch strings.ToLower(key) {
			case "authorization":
				// metadata keys в gRPC должны быть lowercase
				return "authorization", true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := grpcEndpointForGateway(grpcAddr)
	if err := auth_pb.RegisterAuthActionServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening", "addr", httpAddr, "grpc_endpoint", endpoint)
	return http.ListenAndServe(httpAddr, r)
}

func serveSwaggerJSONWithAuth(w http.ResponseWriter, swaggerPath string) {
	b, err := os.ReadFile(swaggerPath)
	if err != nil {
		http.Error(w, "failed to read swagger", http.StatusInternalServerError)
		return
	}

	var doc map[string]any
	if err := json.Unmarshal(b, &doc); err != nil {
		// если вдруг файл битый — отдаём как есть
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
		return
	}

	// OpenAPI v2: добавляем Bearer авторизацию, чтобы Swagger UI показал кнопку Authorize
	doc["securityDefinitions"] = map[string]any{
		"BearerAuth": map[string]any{
			"type":        "apiKey",
			"name":        "Authorization",
			"in":          "header",
			"description": "Bearer <token>",
		},
	}
	doc["security"] = []any{
		map[string]any{"BearerAuth": []any{}},
	}

	out, err := json.Marshal(doc)
	if err != nil {
		http.Error(w, "failed to marshal swagger", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}

func grpcEndpointForGateway(grpcAddr string) string {
	a := strings.TrimSpace(grpcAddr)
	if strings.HasPrefix(a, ":") {
		return "localhost" + a
	}
	return a
}
