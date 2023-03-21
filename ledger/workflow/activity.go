package workflow

import (
	"context"

	"encore.app/ledger/repository"
)

type Activities struct {
	LedgerRepo repository.LedgerRepository
}

func (a *Activities) ExpireAuthorization(ctx context.Context, id string) (string, error) {
	transferId, _ := a.LedgerRepo.VoidPendingTransfer(id)
	return transferId, nil
}
