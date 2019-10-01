package txlistener

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/producer"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
)

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	// Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope, msg)
	if err != nil {
		return err
	}

	// Set Topic to Nonce topic
	msg.Topic = utils.KafkaChainTopic(viper.GetString("kafka.topic.decoder"), txctx.Envelope.GetChain().ID())

	// Set key
	Sender := txctx.Envelope.GetFrom().Address()
	msg.Key = sarama.StringEncoder(utils.ToChainAccountKey(txctx.Envelope.GetChain().ID(), Sender))

	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}