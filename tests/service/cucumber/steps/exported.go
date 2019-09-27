package steps

import (
	"context"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
)

// Init initialize handlers
func Init(ctx context.Context) {
	broker.InitSyncProducer(ctx)
}