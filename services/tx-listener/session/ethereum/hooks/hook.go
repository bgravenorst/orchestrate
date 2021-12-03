package hook

import (
	"context"


	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
	"github.com/consensys/orchestrate/services/tx-listener/entities"
)

//go:generate mockgen -source=hook.go -destination=mock/mock.go -package=mock

type Hook interface {
	AfterNewBlock(ctx context.Context, chain *dynamic.Chain, block *entities.FetchedBlock) error
}
