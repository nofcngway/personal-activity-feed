package bootstrap

import (
	"github.com/nofcngway/feed-service/internal/services/feedservice"
	useractionsprocessor "github.com/nofcngway/feed-service/internal/services/processors/user_actions_processor"
)

func InitUserActionsProcessor(feedService *feedservice.Service) *useractionsprocessor.UserActionsProcessor {
	return useractionsprocessor.New(feedService)
}
