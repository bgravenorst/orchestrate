// +build unit

package static_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/testutils"
)

func TestProvider(t *testing.T) {
	p := static.New(&testutils.Message{Conf: "test-conf"})
	msgs := make(chan provider.Message, 1)
	defer close(msgs)

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = p.Provide(ctx, msgs)
		close(done)
	}()

	msg, _ := (<-msgs).(*testutils.Message)
	assert.Equal(t, "test-conf", msg.Conf, "Message should have flowed properly")
	cancel()
	<-done
}