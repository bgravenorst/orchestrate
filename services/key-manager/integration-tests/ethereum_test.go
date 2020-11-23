// +build integration

package integrationtests

import (
	"github.com/consensys/quorum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	"testing"
)

// keyManagerEthereumTestSuite is a test suite for Key Manager Ethereum controller
type keyManagerEthereumTestSuite struct {
	suite.Suite
	baseURL string
	client  client.KeyManagerClient
	env     *IntegrationEnvironment
}

func (s *keyManagerEthereumTestSuite) SetupSuite() {
	conf := client.NewConfig(s.baseURL, nil)
	s.client = client.NewHTTPClient(http.NewClient(http.NewDefaultConfig()), conf)
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Create() {
	ctx := s.env.ctx

	s.T().Run("should create a new account successfully", func(t *testing.T) {
		accountRequest := testutils.FakeCreateETHAccountRequest()

		account, err := s.client.ETHCreateAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.True(t, common.IsHexAddress(account.Address))
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Import() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = ""

		_, err := s.client.ETHImportAccount(ctx, accountRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should create a new account successfully", func(t *testing.T) {
		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"

		account, err := s.client.ETHImportAccount(ctx, accountRequest)

		assert.NoError(t, err)
		assert.Equal(t, "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5", account.Address)
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_Sign() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := &keymanager.PayloadRequest{
			Data: "",
		}

		_, err := s.client.ETHSign(ctx, "0xaddress", signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should sign payload successfully", func(t *testing.T) {
		expectedAddress := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5"

		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		account, err := s.client.ETHImportAccount(ctx, accountRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, account.Address)

		signRequest := &keymanager.PayloadRequest{
			Data:      "my data to sign",
			Namespace: "_",
		}
		signature, err := s.client.ETHSign(ctx, expectedAddress, signRequest)
		assert.NoError(t, err)
		assert.Equal(t, "0x9a0a890215ea6e79d06f9665297996ab967db117f36c2090d6d6ead5a2d32d5265bc4bc766b5a833cb58b3319e44e952487559b9b939cb5268c0409398214c8b00", signature)
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_SignTransaction() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := testutils.FakeSignETHTransactionRequest()
		signRequest.ChainID = ""

		_, err := s.client.ETHSignTransaction(ctx, "0xaddress", signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should sign transaction successfully", func(t *testing.T) {
		expectedAddress := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5"

		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		account, err := s.client.ETHImportAccount(ctx, accountRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, account.Address)

		signRequest := testutils.FakeSignETHTransactionRequest()
		signature, err := s.client.ETHSignTransaction(ctx, expectedAddress, signRequest)
		assert.NoError(t, err)
		assert.Equal(t, "0x3dcedc00acdd28aab04c2b352608fc6a3cbb17a82935c9168f434ee6d85ddbdd6c75f3299b37977796c019825c9ef49626fd805daa46efc495c5abb2e836446a01", signature)
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_SignQuorumPrivateTransaction() {
	ctx := s.env.ctx

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		signRequest.GasLimit = 0

		_, err := s.client.ETHSignQuorumPrivateTransaction(ctx, "0xaddress", signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should sign transaction successfully", func(t *testing.T) {
		expectedAddress := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5"

		accountRequest := testutils.FakeImportETHAccountRequest()
		accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
		account, err := s.client.ETHImportAccount(ctx, accountRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedAddress, account.Address)

		signRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		signature, err := s.client.ETHSignQuorumPrivateTransaction(ctx, expectedAddress, signRequest)
		assert.NoError(t, err)
		assert.Equal(t, "0x956f2768faa93fbee46bac2fa357c6966401ba7f1b1041125aeb6a4a707088dd6b4d7f697cb456ac2fe58f984da18c03277d53fb67fd429f8f5ba056f5f858ba01", signature)
	})
}

func (s *keyManagerEthereumTestSuite) TestKeyManager_Ethereum_SignEEATransaction() {
	ctx := s.env.ctx

	accountRequest := testutils.FakeImportETHAccountRequest()
	accountRequest.PrivateKey = "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249"
	expectedAddress := "0xeca84382E0f1dDdE22EedCd0D803442972EC7BE5"
	_, _ = s.client.ETHImportAccount(ctx, accountRequest)

	s.T().Run("should fail if payload is invalid", func(t *testing.T) {
		signRequest := testutils.FakeSignEEATransactionRequest()
		signRequest.PrivateFor = []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="}
		signRequest.PrivacyGroupID = "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="

		_, err := s.client.ETHSignEEATransaction(ctx, expectedAddress, signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should fail if privateFor is not base64", func(t *testing.T) {
		signRequest := testutils.FakeSignEEATransactionRequest()
		signRequest.PrivateFrom = "invalid base 64"

		_, err := s.client.ETHSignEEATransaction(ctx, expectedAddress, signRequest)

		assert.True(t, errors.IsInvalidFormatError(err))
	})

	s.T().Run("should sign transaction successfully with privateFor", func(t *testing.T) {
		signRequest := testutils.FakeSignEEATransactionRequest()
		signRequest.PrivateFor = []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="}
		signRequest.PrivacyGroupID = ""

		signature, err := s.client.ETHSignEEATransaction(ctx, expectedAddress, signRequest)

		assert.NoError(t, err)
		assert.Equal(t, "0x4d3803e2633c35f35a18c715aa27548a1ac9e7bc8234aa56322f4a04a71cfeb45f4671bc0f81ae84f11c2397d1d9cf171464f29bfdd381b55deecaa8340c33fb00", signature)
	})

	s.T().Run("should sign transaction successfully with privacyGroupID", func(t *testing.T) {
		signRequest := testutils.FakeSignEEATransactionRequest()
		signRequest.PrivateFor = []string{}
		signRequest.PrivacyGroupID = "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="

		signature, err := s.client.ETHSignEEATransaction(ctx, expectedAddress, signRequest)

		assert.NoError(t, err)
		assert.Equal(t, "0xd1a70ac2784a858394a13189cdea90876d3f34e1d850426a2cd8f71b5c4333d85b55970e274736011d0e5c55c5a8a310a71a8bff424ce1523d173a395644fb7601", signature)
	})
}