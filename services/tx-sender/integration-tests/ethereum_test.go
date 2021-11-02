// +build integration

package integrationtests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	http2 "net/http"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/consensys/orchestrate/pkg/encoding/json"
	encoding "github.com/consensys/orchestrate/pkg/encoding/proto"
	quorumkeymanager "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	utils2 "github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"github.com/consensys/orchestrate/pkg/types/tx"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"gopkg.in/h2non/gock.v1"
)

const (
	waitForEnvelopeTimeOut = 5 * time.Second
)

type txSenderEthereumTestSuite struct {
	suite.Suite
	env *IntegrationEnvironment
}

func (s *txSenderEthereumTestSuite) TestTxSender_Ethereum_Public() {
	signedRawTx := "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
	txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"

	s.T().Run("should sign and send public ethereum transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedRawTx)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", signedRawTx)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		require.NoError(s.T(), err)

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should craft, sign and send public ethereum transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedRawTx)

		feeHistory := testutils.FakeFeeHistory(new(big.Int).SetUint64(100000))
		feeHistoryResult, _ := json.Marshal(feeHistory)
		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_feeHistory")).
			Reply(200).BodyString(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":%s}", feeHistoryResult))

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_estimateGas")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x5208\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", signedRawTx)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "2500100000")).
			Reply(200).JSON(&api.JobResponse{})

		envelope.GasPrice = nil
		envelope.Gas = nil
		envelope.Nonce = nil
		envelope = envelope.SetPriority(utils.PriorityVeryHigh)
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should craft, sign, fail to send, and resend public ethereum transaction successfully(tx-legacy)", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedRawTx)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_gasPrice")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x989680\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_estimateGas")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x5208\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", signedRawTx)).
			Reply(429).BodyString("")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusResending, "", "6000000", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", signedRawTx)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		envelope.GasPrice = nil
		envelope.Gas = nil
		envelope.Nonce = nil
		envelope.TransactionType = string(entities.LegacyTxType)
		envelope = envelope.SetPriority(utils.PriorityVeryLow)
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should sign and send a public onetimekey ethereum transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope().SetTransactionType(string(entities.LegacyTxType))
		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		// IMPORTANT: As we cannot infer txHash before hand, status will be updated to WARNING
		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusWarning, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should send envelope to tx-recover if key-manager fails to sign", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: "not_allowed_tenant",
			},
		})

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
			waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})

	s.T().Run("should send envelope to tx-recover if transaction-scheduler fails to update", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedRawTx)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(422).JSON(httputil.ErrorResponse{
			Message: "cannot update status",
			Code:    666,
		})

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		envelope.Nonce = nil
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
			waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})

	s.T().Run("should send envelope to tx-recover if chain-proxy response with an error", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedRawTx)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"error\":\"invalid_raw\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
			waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})
}

func (s *txSenderEthereumTestSuite) TestTxSender_Ethereum_Raw_Public() {
	raw := "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
	txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"

	s.T().Run("should send ethereum raw transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		_ = envelope.SetRawString(raw)
		_ = envelope.SetTxHashString(txHash)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", raw)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should resend ethereum raw transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		_ = envelope.SetRawString(raw)
		_ = envelope.SetTxHashString(txHash)

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", raw)).
			Reply(429).BodyString("")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusResending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", raw)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should send envelope to tx-recover if send ethereum raw transaction fails", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_RAW_TX)
		_ = envelope.SetRawString(raw)
		_ = envelope.SetTxHashString(txHash)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction", raw)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"invalid_raw\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
			waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})
}

func (s *txSenderEthereumTestSuite) TestTxSender_Ethereum_EEA() {
	txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"

	s.T().Run("should sign and send a EEA transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		signedRawTx := "0xf8be8080808083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1ea0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564"

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-eea-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedRawTx)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "priv_getEeaTransactionCount", envelope.From.String(),
				envelope.PrivateFrom, envelope.PrivateFor)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "priv_distributeRawTransaction", signedRawTx)).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusStored, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		envelope.Nonce = nil
		envelope.Gas = nil
		envelope.GasPrice = nil
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should sign and send EEA transaction with oneTimeKey successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "priv_distributeRawTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusStored, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should send a envelope to tx-recover if we fail to send EEA transaction", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_ORION_EEA_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "priv_distributeRawTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"invalid_raw\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
			waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})
}

func (s *txSenderEthereumTestSuite) TestTxSender_Tessera_Marking() {
	signedTxRaw := "0xf851808398968082520880839896808026a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
	txHash := "0x226d79b217b5ebfeddd08662f3ae1bb1b2cb339d50bbcb708b53ad5f4c71c5ea"

	s.T().Run("should sign and send Tessera marking transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		envelope.Nonce = nil
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-quorum-private-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedTxRaw)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawPrivateTransaction", signedTxRaw, map[string]interface{}{
				"privateFor":  envelope.PrivateFor,
				"privacyFlag": entities.PrivacyFlagSP,
			})).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should sign, fail to send, and resend Tessera marking transaction successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		envelope.Nonce = nil
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-quorum-private-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedTxRaw)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawPrivateTransaction", signedTxRaw, map[string]interface{}{
				"privateFor": envelope.PrivateFor,
				"privacyFlag": entities.PrivacyFlagSP,
			})).
			Reply(429).BodyString("")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusResending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawPrivateTransaction", signedTxRaw, map[string]interface{}{
				"privateFor": envelope.PrivateFor,
				"privacyFlag": entities.PrivacyFlagSP,
			})).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should sign and send Tessera marking transaction with oneTimeKey successfully", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawPrivateTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		// HASH won't match
		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusWarning, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		envelope.Nonce = nil
		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest().EnableTxFromOneTimeKey())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)
	})

	s.T().Run("should send a envelope to tx-recover if we fail to send Tessera marking transaction", func(t *testing.T) {
		defer gock.Off()
		wg := &multierror.Group{}

		envelope := fakeEnvelope()
		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_MARKING_TX)
		_ = envelope.SetTransactionType(string(entities.LegacyTxType))

		gock.New(keyManagerURL).
			Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
			Reply(200).JSON(&types.EthAccountResponse{
			Address: ethcommon.HexToAddress(envelope.GetFromString()),
			Tags: map[string]string{
				quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
			},
		})

		gock.New(keyManagerURL).
			Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-quorum-private-transaction", qkmStoreName, envelope.GetFromString())).
			Reply(200).BodyString(signedTxRaw)

		gock.New(apiURL).
			Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
			AddMatcher(ethCallMatcher(wg, "eth_sendRawPrivateTransaction")).
			Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"error\":\"invalid_raw\"}")

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		gock.New(apiURL).
			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
			Reply(200).JSON(&api.JobResponse{})

		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		err = waitTimeout(wg, waitForEnvelopeTimeOut)
		assert.NoError(t, err)

		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
			waitForEnvelopeTimeOut)
		if err != nil {
			assert.Fail(t, err.Error())
			return
		}

		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
	})
}

// func (s *txSenderEthereumTestSuite) TestTxSender_Tessera_Private() {
// 	data := "0xf8c380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1ea0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564"
// 	enclaveKey := "0x226d79b217b5ebfeddd08662f3ae1bb1b2cb339d50bbcb708b53ad5f4c71c5ea"
// 
// 	s.T().Run("should send a Tessera private transaction successfully", func(t *testing.T) {
// 		defer gock.Off()
// 		wg := &multierror.Group{}
// 
// 		envelope := fakeEnvelope()
// 		envelope.SetDataString(data)
// 		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX)
// 		_ = envelope.SetTransactionType(string(entities.LegacyTxType))
// 
// 		encodedKey := base64.StdEncoding.EncodeToString([]byte(enclaveKey))
// 		gock.New(apiURL).
// 			Post(fmt.Sprintf("/tessera/%s/storeraw", envelope.GetChainUUID())).
// 			Reply(200).JSON(rpc.StoreRawResponse{Key: encodedKey})
// 
// 		gock.New(apiURL).
// 			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
// 			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusStored, "","","")).
// 			Reply(200).JSON(&api.JobResponse{})
// 
// 		envelope.Nonce = nil
// 		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
// 		if err != nil {
// 			assert.Fail(t, err.Error())
// 			return
// 		}
// 
// 		err = waitTimeout(wg, waitForEnvelopeTimeOut)
// 		assert.NoError(t, err)
// 	})
// 
// 	s.T().Run("fail to send a Tessera private transaction", func(t *testing.T) {
// 		defer gock.Off()
// 		wg := &multierror.Group{}
// 
// 		envelope := fakeEnvelope()
// 		envelope.SetDataString(data)
// 		_ = envelope.SetJobType(tx.JobType_ETH_TESSERA_PRIVATE_TX)
// 		_ = envelope.SetTransactionType(string(entities.LegacyTxType))
// 
// 		gock.New(apiURL).
// 			Post(fmt.Sprintf("/tessera/%s/storeraw", envelope.GetChainUUID())).
// 			Reply(200).BodyString("Invalid_Response")
// 
// 		gock.New(apiURL).
// 			Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
// 			AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "","","")).
// 			Reply(200).JSON(&api.JobResponse{})
// 
// 		err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
// 		if err != nil {
// 			assert.Fail(t, err.Error())
// 			return
// 		}
// 
// 		err = waitTimeout(wg, waitForEnvelopeTimeOut)
// 		assert.NoError(t, err)
// 
// 		retrievedEnvelope, err := s.env.consumer.WaitForEnvelope(envelope.GetID(), s.env.srvConfig.RecoverTopic,
// 			waitForEnvelopeTimeOut)
// 		if err != nil {
// 			assert.Fail(t, err.Error())
// 			return
// 		}
// 
// 		assert.Equal(t, envelope.GetID(), retrievedEnvelope.GetID())
// 		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Code)
// 		assert.NotEmpty(t, retrievedEnvelope.GetErrors()[0].Message)
// 	})
// }

func (s *txSenderEthereumTestSuite) TestTxSender_XNonceManager() {
	signedRawTx := "0xf85380839896808252088083989680808216b4a0d35c752d3498e6f5ca1630d264802a992a141ca4b6a3f439d673c75e944e5fb0a05278aaa5fabbeac362c321b54e298dedae2d31471e432c26ea36a8d49cf08f1e"
	txHash := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f3"
	txHash2 := "0x6621fbe1e2848446e38d99bfda159cdd83f555ae0ed7a4f3e1c3c79f7d6d74f4"

	s.T().Run("should increment account nonce on consecutive transaction successfully", func(t *testing.T) {
		envelope := fakeEnvelope()
		envelope.Nonce = nil

		for idx := 0; idx < 3; idx++ {
			wg := &multierror.Group{}
			_ = envelope.SetJobUUID(uuid.Must(uuid.NewV4()).String())

			if idx == 0 {
				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")
			}

			gock.New(keyManagerURL).
				Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
				Reply(200).JSON(&types.EthAccountResponse{
				Address: ethcommon.HexToAddress(envelope.GetFromString()),
				Tags: map[string]string{
					quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
				},
			})

			gock.New(keyManagerURL).
				Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
				Reply(200).BodyString(signedRawTx)

			if idx == 2 {
				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash2 + "\"}")
			} else {
				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")
			}

			gock.New(apiURL).
				Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
				AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, fmt.Sprintf("%d", idx), "", "")).
				Reply(200).JSON(&api.JobResponse{})

			// Warning because txHash does not match
			if idx == 2 {
				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusWarning, fmt.Sprintf("%d", idx), "", "")).
					Reply(200).JSON(&api.JobResponse{})
			}

			err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			err = waitTimeout(wg, waitForEnvelopeTimeOut)
			assert.NoError(t, err)
			gock.Off()
		}

		nonce, _, _ := s.env.ns.GetLastSent(envelope.PartitionKey())
		assert.Equal(t, uint64(2), nonce)
	})

	s.T().Run("should re-fetch nonce on nonce too low errors", func(t *testing.T) {
		envelope := fakeEnvelope()
		envelope.Nonce = nil

		for idx := 0; idx < 3; idx++ {
			wg := &multierror.Group{}
			_ = envelope.SetJobUUID(uuid.Must(uuid.NewV4()).String())

			if idx == 0 {
				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

				gock.New(keyManagerURL).
					Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
					Reply(200).JSON(&types.EthAccountResponse{
					Address: ethcommon.HexToAddress(envelope.GetFromString()),
					Tags: map[string]string{
						quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
					},
				})

				gock.New(keyManagerURL).
					Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
					Reply(200).BodyString(signedRawTx)

				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, fmt.Sprintf("%d", idx), "", "")).
					Reply(200).JSON(&api.JobResponse{})
			}

			if idx == 1 {
				gock.New(keyManagerURL).
					Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).Times(2).
					Reply(200).JSON(&types.EthAccountResponse{
					Address: ethcommon.HexToAddress(envelope.GetFromString()),
					Tags: map[string]string{
						quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
					},
				})

				gock.New(keyManagerURL).
					Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).Times(2).
					Reply(200).BodyString(signedRawTx)

				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "1", "", "")).
					Reply(200).JSON(&api.JobResponse{})

				resp := utils2.JSONRpcMessage{Error: &utils2.JSONError{Code: 100, Message: "nonce too low"}}
				bresp, _ := json.Marshal(resp)
				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
					Reply(200).BodyString(string(bresp))

				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x2\"}")

				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusRecovering, "", "", "")).
					Reply(200).JSON(&api.JobResponse{})

				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash + "\"}")

				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "2", "", "")).
					Reply(200).JSON(&api.JobResponse{})
			}

			if idx > 1 {
				gock.New(keyManagerURL).
					Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
					Reply(200).JSON(&types.EthAccountResponse{
					Address: ethcommon.HexToAddress(envelope.GetFromString()),
					Tags: map[string]string{
						quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
					},
				})

				gock.New(keyManagerURL).
					Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
					Reply(200).BodyString(signedRawTx)

				gock.New(apiURL).
					Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
					AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
					Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"" + txHash2 + "\"}")

				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "3", "", "")).
					Reply(200).JSON(&api.JobResponse{})

				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusWarning, "3", "", "")).
					Reply(200).JSON(&api.JobResponse{})
			}

			err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			err = waitTimeout(wg, waitForEnvelopeTimeOut)
			assert.NoError(t, err)
			gock.Off()
		}

		nonce, _, _ := s.env.ns.GetLastSent(envelope.PartitionKey())
		assert.Equal(t, uint64(3), nonce)
	})

	s.T().Run("should retry on nonce too low errors till max recover", func(t *testing.T) {
		envelope := fakeEnvelope()
		envelope.Nonce = nil

		for idx := 0; idx <= maxRecoveryDefault; idx++ {
			wg := &multierror.Group{}

			gock.New(apiURL).
				Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
				AddMatcher(ethCallMatcher(wg, "eth_getTransactionCount", envelope.From.String(), "pending")).
				Reply(200).BodyString("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x0\"}")

			gock.New(keyManagerURL).
				Get(fmt.Sprintf("/stores/%s/ethereum/%s", qkmStoreName, envelope.GetFromString())).
				Reply(200).JSON(&types.EthAccountResponse{
				Address: ethcommon.HexToAddress(envelope.GetFromString()),
				Tags: map[string]string{
					quorumkeymanager.TagIDAllowedTenants: multitenancy.DefaultTenant,
				},
			})

			gock.New(keyManagerURL).
				Post(fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", qkmStoreName, envelope.GetFromString())).
				Reply(200).BodyString(signedRawTx)

			gock.New(apiURL).
				Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
				AddMatcher(txStatusUpdateMatcher(wg, entities.StatusPending, "0", "", "")).
				Reply(200).JSON(&api.JobResponse{})

			resp := utils2.JSONRpcMessage{Error: &utils2.JSONError{Code: 100, Message: "nonce too low"}}
			bresp, _ := json.Marshal(resp)
			gock.New(apiURL).
				Post(fmt.Sprintf("/proxy/chains/%s", envelope.GetChainUUID())).
				AddMatcher(ethCallMatcher(wg, "eth_sendRawTransaction")).
				Reply(200).BodyString(string(bresp))

			if idx < maxRecoveryDefault {
				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusRecovering, "", "", "")).
					Reply(200).JSON(&api.JobResponse{})
			} else {
				gock.New(apiURL).
					Patch(fmt.Sprintf("/jobs/%s", envelope.GetJobUUID())).
					AddMatcher(txStatusUpdateMatcher(wg, entities.StatusFailed, "", "", "")).
					Reply(200).JSON(&api.JobResponse{})
			}

			if idx == 0 {
				err := s.sendEnvelope(envelope.TxEnvelopeAsRequest())
				if err != nil {
					assert.Fail(t, err.Error())
					return
				}
			}

			err := waitTimeout(wg, waitForEnvelopeTimeOut)
			assert.NoError(t, err)

			gock.Off()
		}

		_, ok, _ := s.env.ns.GetLastSent(envelope.PartitionKey())
		assert.False(t, ok)
	})
}

func (s *txSenderEthereumTestSuite) TestTxSender_ZHealthCheck() {
	type healthRes struct {
		API   string `json:"api,omitempty"`
		Kafka string `json:"kafka,omitempty"`
		Redis string `json:"redis,omitempty"`
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	ctx := s.env.ctx
	s.T().Run("should retrieve positive health check over service dependencies", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(apiMetricsURL).Get("/live").Reply(200)
		defer gock.Off()

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		assert.Equal(s.T(), 200, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "OK", status.API)
		assert.Equal(s.T(), "OK", status.Kafka)
		assert.Equal(s.T(), "OK", status.Redis)
	})

	s.T().Run("should retrieve a negative health check over kafka service", func(t *testing.T) {
		req, err := http2.NewRequest("GET", fmt.Sprintf("%s/ready?full=1", s.env.metricsURL), nil)
		assert.NoError(s.T(), err)

		gock.New(apiMetricsURL).Get("/live").Reply(200)
		defer gock.Off()

		// Kill Kafka on first call so data is added in DB and status is CREATED but does not get updated to STARTED
		err = s.env.client.Stop(ctx, kafkaContainerID)
		assert.NoError(t, err)

		resp, err := httpClient.Do(req)
		if err != nil {
			assert.Fail(s.T(), err.Error())
			return
		}

		err = s.env.client.StartServiceAndWait(ctx, kafkaContainerID, 10*time.Second)
		assert.NoError(t, err)

		assert.Equal(s.T(), 503, resp.StatusCode)
		status := healthRes{}
		err = json.UnmarshalBody(resp.Body, &status)
		assert.NoError(s.T(), err)
		assert.NotEqual(s.T(), "OK", status.Kafka)
		assert.Equal(s.T(), "OK", status.API)
		assert.Equal(s.T(), "OK", status.Redis)
	})
}

func fakeEnvelope() *tx.Envelope {
	scheduleUUID := uuid.Must(uuid.NewV4()).String()
	jobUUID := uuid.Must(uuid.NewV4()).String()
	chainUUID := uuid.Must(uuid.NewV4()).String()

	envelope := tx.NewEnvelope()
	_ = envelope.SetID(scheduleUUID)
	_ = envelope.SetJobUUID(jobUUID)
	_ = envelope.SetJobType(tx.JobType_ETH_TX)
	_ = envelope.SetNonce(0)
	_ = envelope.SetFromString(ethcommon.HexToAddress(utils.RandHexString(12)).String())
	_ = envelope.SetGas(21000)
	_ = envelope.SetGasPriceString("10000000")
	_ = envelope.SetValueString("10000000")
	_ = envelope.SetDataString("0x")
	_ = envelope.SetChainIDString("2888")
	_ = envelope.SetChainUUID(chainUUID)
	_ = envelope.SetHeadersValue(utils.TenantIDMetadata, "tenantID")
	_ = envelope.SetPrivateFrom("A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=")
	_ = envelope.SetPrivateFor([]string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="})

	return envelope
}

func (s *txSenderEthereumTestSuite) sendEnvelope(protoMessage proto.Message) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = s.env.srvConfig.SenderTopic

	b, err := encoding.Marshal(protoMessage)
	if err != nil {
		return err
	}
	msg.Value = sarama.ByteEncoder(b)

	_, _, err = s.env.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func waitTimeout(wg *multierror.Group, duration time.Duration) error {
	c := make(chan bool, 1)
	var err error
	go func() {
		defer close(c)
		err = wg.Wait().ErrorOrNil()
	}()

	select {
	case <-c:
		return err
	case <-time.After(duration):
		return fmt.Errorf("timeout after %s", duration.String())
	}
}

func txStatusUpdateMatcher(wg *multierror.Group, status entities.JobStatus, nonce string, gasPrice string, gasFeeCap string) gock.MatchFunc {
	cerr := make(chan error, 1)
	wg.Go(func() error {
		return <-cerr
	})

	return func(rw *http2.Request, _ *gock.Request) (bool, error) {
		defer func() {
			cerr <- nil
		}()

		body, _ := ioutil.ReadAll(rw.Body)
		rw.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		req := api.UpdateJobRequest{}
		if err := json.Unmarshal(body, &req); err != nil {
			cerr <- err
			return false, err
		}

		if status != "" && req.Status != status {
			err := fmt.Errorf("invalid status, got %s, expected %s", req.Status, status)
			cerr <- err
			return false, err
		}
		if nonce != "" && req.Transaction.Nonce != nonce {
			err := fmt.Errorf("invalid nonce, got %s, expected %s", req.Transaction.Nonce, nonce)
			cerr <- err
			return false, err
		}
		if gasPrice != "" && req.Transaction.GasPrice != gasPrice {
			err := fmt.Errorf("invalid gasPrice, got %s, expected %s", req.Transaction.GasPrice, gasPrice)
			cerr <- err
			return false, err
		}
		if gasFeeCap != "" && req.Transaction.GasFeeCap != gasFeeCap {
			err := fmt.Errorf("invalid gasFeeCap, got %s, expected %s", req.Transaction.GasFeeCap, gasFeeCap)
			cerr <- err
			return false, err
		}

		return true, nil
	}
}

func ethCallMatcher(wg *multierror.Group, method string, args ...interface{}) gock.MatchFunc {
	cerr := make(chan error, 1)
	wg.Go(func() error {
		select {
		case err := <-cerr:
			return err
		case <-time.After(waitForEnvelopeTimeOut):
			return fmt.Errorf("timeout after %s", waitForEnvelopeTimeOut.String())
		}
	})

	return func(rw *http2.Request, grw *gock.Request) (bool, error) {
		defer func() {
			cerr <- nil
		}()

		body, _ := ioutil.ReadAll(rw.Body)
		rw.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		req := utils2.JSONRpcMessage{}
		if err := json.Unmarshal(body, &req); err != nil {
			cerr <- err
			return false, err
		}

		if req.Method != method {
			err := fmt.Errorf("invalid method, got %s, expected %s", req.Method, method)
			cerr <- err
			return false, err
		}

		if len(args) > 0 {
			params, _ := json.Marshal(args)
			if strings.ToLower(string(req.Params)) != strings.ToLower(string(params)) {
				err := fmt.Errorf("invalid params, got %s, expected %s", req.Params, params)
				cerr <- err
				return false, err
			}
		}

		return true, nil
	}
}
