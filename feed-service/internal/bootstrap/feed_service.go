package bootstrap

import (
	"github.com/nofcngway/feed-service/internal/services/feedservice"
	"github.com/nofcngway/feed-service/internal/storage/pgstorage"
)

func InitFeedService(storage *pgstorage.PGStorage) *feedservice.Service {
	return feedservice.New(storage)
}


