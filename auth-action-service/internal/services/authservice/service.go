package authservice

import "time"

type Service struct {
	storage    UserStorage
	sessions   SessionStore
	producer   Producer
	sessionTTL time.Duration
}

func New(storage UserStorage, sessionsStore SessionStore, producer Producer, sessionTTL time.Duration) *Service {
	if sessionTTL <= 0 {
		sessionTTL = time.Hour
	}
	return &Service{
		storage:    storage,
		sessions:   sessionsStore,
		producer:   producer,
		sessionTTL: sessionTTL,
	}
}


