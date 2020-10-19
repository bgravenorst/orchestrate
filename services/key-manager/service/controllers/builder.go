package controllers

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/use-cases/ethereum"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

//go:generate swag init --dir . --generalInfo builder.go --output ../../../../public/swagger-specs/services/key-manager
//go:generate rm ../../../../public/swagger-specs/services/key-manager/docs.go ../../../../public/swagger-specs/services/key-manager/swagger.yaml

// @title Key Management API
// @version 2.0
// @description PegaSys Orchestrate Key Management. Enables fine-grained management of cryptographic keys.

// @contact.name Contact PegaSys Orchestrate
// @contact.url https://pegasys.tech/contact/
// @contact.email support@pegasys.tech

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

type Builder struct {
	ethereumCtrl *EthereumController
}

func NewBuilder(ucs ethereum.UseCases) *Builder {
	return &Builder{
		ethereumCtrl: NewEthereumController(ucs),
	}
}

func (b *Builder) Build(ctx context.Context, _ string, configuration interface{}, respModifier func(response *http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Signer)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	b.ethereumCtrl.Append(router)

	return router, nil
}