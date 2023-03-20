package repository

import tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"

type LedgerRepository interface {
	Init()
	CreateAccount(id string) error
	GetAccount(id string) (*tb_types.Account, error)
	CreatePendingTransfer(id string, amount int) (string, error)
	PostPendingTransfer(pendingId string) (string, error)
	VoidPendingTransfer(pendingId string) (string, error)
}
