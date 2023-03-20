package workflow

import (
	"context"
	"fmt"

	"encore.app/ledger/service"
)

type Activities struct {
	LedgerService service.LedgerService
}

func ComposeGreeting(ctx context.Context, name string) (string, error) {
	greeting := fmt.Sprintf("Hello %s!", name)
	return greeting, nil
}

func (a *Activities) ExpireAuthorization(ctx context.Context, id string) (string, error) {
	res, _ := a.LedgerService.VoidPendingPayment(id)
	return res, nil
}
