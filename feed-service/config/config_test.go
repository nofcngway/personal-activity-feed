package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_OK(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`
grpc:
  addr: ":50052"
http:
  addr: ":8083"
  swagger_path: "./swagger.json"
database:
  host: "localhost"
  port: 55432
  username: "admin"
  password: "admin"
  name: "activity"
  ssl_mode: "disable"
kafka:
  brokers: ["localhost:9092"]
  topic_name: "user-actions"
  group_id: "feed-service-group"
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.GRPC.Addr != ":50052" || cfg.HTTP.Addr != ":8083" || cfg.Database.Port != 55432 || cfg.Kafka.TopicName != "user-actions" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := LoadConfig("/no/such/file.yaml")
	if err == nil {
		t.Fatalf("expected error")
	}
}
