package wallet

import (
	"context"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/multi-vault/keystore"
)

func makeGeneratorContext(i int) *engine.TxContext {
	txctx := engine.NewTxContext()
	txctx.Reset()
	txctx.Logger = log.NewEntry(log.StandardLogger())
	switch i % 2 {
	case 0:
		txctx.Envelope.Chain = chain.FromInt(0)
		txctx.Set("errors", 0)
	case 1:
		txctx.Envelope.Chain = chain.FromInt(10)
		txctx.Set("errors", 0)
	}
	return txctx
}

type GeneratorSuite struct {
	testutils.HandlerTestSuite
}

func (s *GeneratorSuite) SetupSuite() {
	// The default keystore uses a mocked secret store
	keystore.Init(context.Background())
	s.Handler = Generator(keystore.GlobalKeyStore())
}

func (s *GeneratorSuite) TestGenerator() {
	rounds := 100
	txctxs := []*engine.TxContext{}

	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeGeneratorContext(i))
	}

	s.Handle(txctxs)

	for _, txctx := range txctxs {
		// Handle contexts

		assert.Len(s.T(), txctx.Envelope.Errors, txctx.Get("errors").(int), "Expected right count of errors")
		assert.NotEqual(s.T(), txctx.Envelope.Sender(), ethcommon.Address{}, "Expected new address to be set")
	}
}

func TestFaucet(t *testing.T) {
	suite.Run(t, new(GeneratorSuite))
}