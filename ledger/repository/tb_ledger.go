package repository

import (
	"strconv"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/rlog"
	tb "github.com/tigerbeetledb/tigerbeetle-go"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

const (
	// Change this value according to your cluster setup
	address     = "127.0.0.1:3000"
	clusterId   = 0
	concurrency = 1
)

type TBLedgerRepository struct {
}

var (
	db tb.Client
)

func (r *TBLedgerRepository) Init() {
	client, err := tb.NewClient(clusterId, []string{address}, concurrency)
	if err != nil {
		rlog.Error("Connection to TB Failed", err)
		panic("Connection to TB Failed")
	}
	rlog.Info("Connected to TB")
	db = client
}

func (r *TBLedgerRepository) CreateAccount(id string) error {
	ctx := rlog.With("id", id)
	res, err := db.CreateAccounts(createAccountArg(id, 1))

	if err != nil {
		ctx.Error("Error creating account", "err", err)
	}

	for _, err := range res {
		ctx.Error("Error creating account", "err", err.Result.String())
		return errs.B().Msg("Error creating account").Err()
	}

	return nil
}

func (r *TBLedgerRepository) GetAccount(id string) (*tb_types.Account, error) {
	ids := []tb_types.Uint128{uint128(id)}
	accounts, err := db.LookupAccounts(ids)
	account := accounts[0]

	return &account, err
}

func (r *TBLedgerRepository) CreatePendingTransfer(id string, amount int) (string, error) {
	// temporarily
	transferId := generateId()
	transfer := tb_types.Transfer{
		ID:        uint128(transferId),
		PendingID: tb_types.Uint128{},
		// Deduct Customer
		DebitAccountID: uint128(id),
		// Default to bank
		CreditAccountID: uint128("1"),
		UserData:        tb_types.Uint128{},
		Reserved:        tb_types.Uint128{},
		Timeout:         0,
		Ledger:          1,
		Code:            1,
		// Pending Transfer
		Flags:     tb_types.TransferFlags{Pending: true}.ToUint16(),
		Amount:    uint64(amount),
		Timestamp: 0,
	}

	transfersRes, err := db.CreateTransfers([]tb_types.Transfer{transfer})

	for _, err := range transfersRes {
		rlog.Error("Batch transfer at %d failed to create: %s", err.Index, err.Result)
	}

	return transferId, err
}

func (r *TBLedgerRepository) PostPendingTransfer(pendingId string) (string, error) {
	transferId := generateId()
	transfer := tb_types.Transfer{
		ID:        uint128(transferId),
		PendingID: uint128(pendingId),
		Flags:     tb_types.TransferFlags{PostPendingTransfer: true}.ToUint16(),
		Timestamp: 0,
	}

	transfersRes, err := db.CreateTransfers([]tb_types.Transfer{transfer})

	for _, err := range transfersRes {
		rlog.Error("Batch transfer at %d failed to create: %s", err.Index, err.Result)
	}

	return transferId, err
}

func (r *TBLedgerRepository) VoidPendingTransfer(pendingId string) (string, error) {
	transferId := generateId()
	transfer := tb_types.Transfer{
		ID:        uint128(transferId),
		PendingID: uint128(pendingId),
		Flags:     tb_types.TransferFlags{VoidPendingTransfer: true}.ToUint16(),
		Timestamp: 0,
	}

	transfersRes, err := db.CreateTransfers([]tb_types.Transfer{transfer})

	for _, err := range transfersRes {
		rlog.Error("Batch transfer at %d failed to create: %s", err.Index, err.Result)
	}

	return transferId, err
}

func uint128(value string) tb_types.Uint128 {
	res, _ := tb_types.HexStringToUint128(value)
	return res
}

func createAccountArg(id string, ledger uint32) []tb_types.Account {
	return []tb_types.Account{{
		ID:             uint128(id),
		UserData:       tb_types.Uint128{},
		Reserved:       [48]uint8{},
		Ledger:         ledger,
		Code:           718,
		Flags:          0,
		DebitsPending:  0,
		DebitsPosted:   0,
		CreditsPending: 0,
		CreditsPosted:  0,
		Timestamp:      0,
	}}
}

func generateId() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
