package authservice

import kafkaproducer "github.com/nofcngway/auth-action-service/internal/kafka/producer"

// compile-time check: наш kafka producer удовлетворяет интерфейсу Producer
var _ Producer = (*kafkaproducer.Producer)(nil)


