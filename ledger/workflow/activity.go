package workflow

import (
	"context"

	"encore.app/ledger/service"
)

type Activities struct {
	LedgerService service.LedgerService
}

func (a *Activities) ExpireAuthorization(ctx context.Context, id string) (string, error) {
	res, _ := a.LedgerService.VoidPendingPayment(id)
	return res, nil
}
