package bootstrap

import (
	server "github.com/nofcngway/feed-service/internal/api/feed_service_api"
	"github.com/nofcngway/feed-service/internal/services/feedservice"
)

func InitFeedAPI(feedService *feedservice.Service) *server.FeedServiceAPI {
	return server.NewFeedServiceAPI(feedService)
}
