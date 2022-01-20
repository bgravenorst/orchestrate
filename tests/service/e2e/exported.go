package e2e

import (
	"context"
	"fmt"
	"sync"

	"encoding/json"

	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/tests/service/e2e/cucumber/alias"
	utils3 "github.com/consensys/orchestrate/tests/utils"

	loader "github.com/consensys/orchestrate/handlers/loader/sarama"
	"github.com/consensys/orchestrate/handlers/offset"
	broker "github.com/consensys/orchestrate/pkg/broker/sarama"
	"github.com/consensys/orchestrate/pkg/engine"
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/tests/handlers"
	"github.com/consensys/orchestrate/tests/handlers/consumer"
	"github.com/consensys/orchestrate/tests/handlers/dispatcher"
	"github.com/consensys/orchestrate/tests/service/e2e/cucumber"
	utils2 "github.com/consensys/orchestrate/tests/service/e2e/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Start starts application
func Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	log.FromContext(ctx).Info("starting execution...")

	var gerr error
	// Create context for application
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cctx = multitenancy.WithUserInfo(cctx, multitenancy.DefaultUser())
	var testData *utils3.TestData
	rawTestData := viper.GetString(e2eDataViperKey)
	logger.WithField("data", rawTestData).Debug("test data")
	err := json.Unmarshal([]byte(rawTestData), &testData)
	if err != nil {
		logger.WithError(err).Error("failed to ")
		return err
	}

	initComponents(cctx, rawTestData)
	registerHandlers()

	err = importTestIdentities(cctx, testData)
	if err != nil {
		return err
	}

	chainUUIDs, err := initTestChains(cctx, testData)
	if err != nil {
		return err
	}

	// Start consuming on every topics of interest
	var topics []string
	for _, viprTopicKey := range utils2.TOPICS {
		topics = append(topics, viper.GetString(viprTopicKey))
	}

	cg := consumer.NewEmbeddingConsumerGroupHandler(engine.GlobalEngine())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.FromContext(cctx).WithField("topics", topics).Info("connecting")

		err := broker.Consume(cctx, topics, cg)
		if err != nil {
			log.FromContext(cctx).WithError(err).Error("error on consumer")
		}

		cancel()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		<-cg.IsReady()
		err := cucumber.Run(cucumber.GlobalOptions())
		if err != nil {
			gerr = err
			log.FromContext(cctx).WithError(err).Error("error on cucumber")
		}

		cancel()
		wg.Done()
	}()

	wg.Wait()

	if err := removeTestChains(ctx, chainUUIDs); err != nil {
		gerr = errors.CombineErrors(gerr, err)
	}

	return gerr
}

func Stop(ctx context.Context) error {
	log.FromContext(ctx).Info("Cucumber: stopping execution...")
	return nil
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(loader.Loader)
	engine.Register(offset.Marker)
	engine.Register(dispatcher.GlobalHandler())
}

func initComponents(ctx context.Context, rawTestData string) {
	utils.InParallel(
		// Initialize Engine
		func() {
			engine.Init(ctx)
		},
		// Initialize ConsumerGroup
		func() {
			broker.InitConsumerGroup(ctx, viper.GetString(broker.ConsumerGroupNameViperKey))
		},
		// Initialize Handlers
		func() {
			// Prepare topics map for dispatcher
			topics := make(map[string]string)
			for topic, viprTopicKey := range utils2.TOPICS {
				topics[viper.GetString(viprTopicKey)] = topic
			}
			dispatcher.SetKeyOfFuncs(
				dispatcher.LongKeyOf(topics),
				dispatcher.LabelKey(topics),
				dispatcher.ShortKeyOf(topics),
			)
			handlers.Init(ctx)
		},
		// Initialize logger
		func() {
			cfg := log.NewConfig(viper.GetViper())
			// Create and configure logger
			logger := logrus.StandardLogger()
			_ = log.ConfigureLogger(cfg, logger)
		},
		// Initialize cucumber handlers
		func() {
			cucumber.Init(ctx, rawTestData)
		},
	)
}

// We import account define at Global Aliases
func importTestIdentities(ctx context.Context, testData *utils3.TestData) error {
	logger := log.FromContext(ctx)
	orchestrateClient := client.GlobalClient()

	nodes := append(testData.Nodes.Besu, testData.Nodes.Quorum...)
	nodes = append(nodes, testData.Nodes.Geth...)
	for idx := range nodes {
		node := nodes[idx]
		for _, privKey := range node.FundedPrivateKeys {
			resp, err := orchestrateClient.ImportAccount(ctx, &api.ImportAccountRequest{
				PrivateKey: privKey,
			})

			if err != nil {
				if errors.IsAlreadyExistsError(err) || errors.IsConflictedError(err) {
					logger.WithError(err).WithField("priv_key", privKey.String()).Warn("imported account is duplicated")
					continue
				}

				logger.WithError(err).WithField("priv_key", privKey.String()).Error("failed to import account")
				return err
			}

			logger.WithField("address", resp.Address.String()).Info("account imported successfully")
		}
	}

	return nil
}

func initTestChains(ctx context.Context, testData *utils3.TestData) (map[string]string, error) {
	aliases := alias.GlobalAliasRegistry()
	logger := log.FromContext(ctx)
	orchestrateClient := client.GlobalClient()
	ec := rpc.GlobalClient()
	proxyHost := viper.GetString(client.URLViperKey)

	reqs := map[string]*api.RegisterChainRequest{}
	for idx := range testData.Nodes.Besu {
		node := testData.Nodes.Besu[idx]
		reqs[fmt.Sprintf("besu%d", idx)] = &api.RegisterChainRequest{
			URLs: node.URLs,
			Name: fmt.Sprintf("besu-%s", utils.RandString(5)),
		}
	}

	for idx := range testData.Nodes.Geth {
		node := testData.Nodes.Geth[idx]
		reqs[fmt.Sprintf("geth%d", idx)] = &api.RegisterChainRequest{
			URLs: node.URLs,
			Name: fmt.Sprintf("geth-%s", utils.RandString(5)),
		}
	}

	for idx := range testData.Nodes.Quorum {
		node := testData.Nodes.Quorum[idx]
		if len(node.URLs) == 0 {
			continue
		}
		req := &api.RegisterChainRequest{
			URLs: node.URLs,
			Name: fmt.Sprintf("quorum-%s", utils.RandString(5)),
		}
		if node.PrivateTxManager.URL != "" {
			req.PrivateTxManager = &api.PrivateTxManagerRequest{
				URL:  node.PrivateTxManager.URL,
				Type: entities.TesseraChainType,
			}
		}
		reqs[fmt.Sprintf("quorum%d", idx)] = req
	}

	chainUUIDs := map[string]string{}
	for chainAlias, req := range reqs {
		resp, err := orchestrateClient.RegisterChain(ctx, req)
		if err != nil {
			logger.WithField("name", req.Name).WithError(err).Error("failed to register chain")
			return chainUUIDs, err
		}

		logger.WithField("name", req.Name).WithField("uuid", resp.UUID).WithField("alias", chainAlias).
			Info("chain registered successfully")
		chainUUIDs[req.Name] = resp.UUID

		aliases.Set(resp.UUID, fmt.Sprintf("chain.%s.UUID", chainAlias))
		aliases.Set(resp.Name, fmt.Sprintf("chain.%s.Name", chainAlias))
	}

	for _, chainUUID := range chainUUIDs {
		err := utils3.WaitForProxy(ctx, proxyHost, chainUUID, ec)
		if err != nil {
			logger.WithField("uuid", chainUUID).WithError(err).Error("failed to wait for proxy chain")
			return chainUUIDs, err
		}
	}

	return chainUUIDs, nil
}

func removeTestChains(ctx context.Context, chainUUIDs map[string]string) error {
	orchestrateClient := client.GlobalClient()
	logger := log.FromContext(ctx)
	for chainName, chainUUID := range chainUUIDs {
		err := orchestrateClient.DeleteChain(ctx, chainUUID)
		if err != nil {
			logger.WithField("uuid", chainUUID).WithField("name", chainName).
				WithError(err).Error("failed to remove test chain")
			return err
		}

		logger.WithField("uuid", chainUUID).WithField("name", chainName).
			Info("test chain was removed successfully")
	}

	return nil
}
