package contracts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const searchContractComponent = "use-cases.get-contract-by-code-hash"

type searchContractUseCase struct {
	agent  store.DB
	logger *log.Logger
}

func NewSearchContractUseCase(agent store.DB) usecases.SearchContractUseCase {
	return &searchContractUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(searchContractComponent),
	}
}

// Execute gets a contract from DB
func (uc *searchContractUseCase) Execute(ctx context.Context, codehash, signhash hexutil.Bytes) (*entities.Contract, error) {
	return nil, errors.NotFoundError("Pending to develop")
	// logger := uc.logger.WithContext(ctx)
	// 
	// artifact, err := uc.agent.Artifact().FindOneByABIAndCodeHash()
	// if err != nil {
	// 	return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	// }
	// 
	// contract := &entities.Contract{
	// 	Name:             name,
	// 	Tag:              tag,
	// 	ABI:              artifact.ABI,
	// 	Bytecode:         hexutil.MustDecode(artifact.Bytecode),
	// 	DeployedBytecode: hexutil.MustDecode(artifact.DeployedBytecode),
	// }
	// 
	// contractABI, err := contract.ToABI()
	// if err != nil {
	// 	errMessage := "failed to get contract ABI"
	// 	logger.WithError(err).Error(errMessage)
	// 	return nil, errors.DataCorruptedError(errMessage).ExtendComponent(getMethodSignaturesComponent)
	// }
	// 
	// // nolint
	// for _, method := range contractABI.Methods {
	// 	contract.Methods = append(contract.Methods, entities.Method{
	// 		Signature: method.Sig,
	// 	})
	// }
	// 
	// // nolint
	// for _, event := range contractABI.Events {
	// 	contract.Events = append(contract.Events, entities.Event{
	// 		Signature: event.Sig,
	// 	})
	// }
	// 
	// if contractABI.Constructor.Sig == "" {
	// 	contract.Constructor = entities.Method{Signature: "()"}
	// } else {
	// 	contract.Constructor = entities.Method{Signature: contractABI.Constructor.Sig}
	// }
	// 
	// logger.Debug("contract found successfully")
	// return contract, nil
}
