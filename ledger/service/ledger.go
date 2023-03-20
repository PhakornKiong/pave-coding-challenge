package service

import (
	"fmt"

	"encore.app/ledger/repository"
	"encore.dev/rlog"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type LedgerService struct {
	LedgerRepo repository.LedgerRepository
}

type AccountBalance struct {
	Available uint64
	Reserved  uint64
}
type LedgerAccount struct {
	Id      string
	Balance AccountBalance
}

func (s *LedgerService) CreateAccount(id string) error {
	err := s.LedgerRepo.CreateAccount(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *LedgerService) GetAccountBalance(id string) (LedgerAccount, error) {
	account, _ := s.LedgerRepo.GetAccount(id)
	ledgerAccount := buildLedgerAccount(account)
	return ledgerAccount, nil
}

func (s *LedgerService) AddAccountBalance(customerId string, amount int) (string, error) {
	// default debit user 1 as our own bank
	transferId, err := s.LedgerRepo.CreateTransfer("1", customerId, amount)
	return transferId, err
}

func (s *LedgerService) CreatePayment(customerId string, amount int) (string, error) {
	// default credit user 1 as our own bank
	transferId, err := s.LedgerRepo.CreateTransfer(customerId, "1", amount)
	return transferId, err
}

func (s *LedgerService) AuthorizePayment(customerId string, amount int) (string, error) {
	transferId, err := s.LedgerRepo.CreatePendingTransfer(customerId, amount)
	return transferId, err
}

func (s *LedgerService) ReleasePayment(id string) (string, error) {
	transferId, _ := s.LedgerRepo.PostPendingTransfer(id)
	return transferId, nil
}

func (s *LedgerService) VoidPendingPayment(id string) (string, error) {
	transferId, _ := s.LedgerRepo.VoidPendingTransfer(id)
	return transferId, nil
}

// Assume Debit is asset and Credit is liability
// Account 1 is the bank, others are customers
// So balance is customer POV is technically credit
// TODO: Handle overflow
func buildLedgerAccount(account *tb_types.Account) LedgerAccount {
	creditPosted := account.CreditsPosted
	debitPosted, debitPending := account.DebitsPosted, account.DebitsPending
	rlog.Info(fmt.Sprint(debitPosted, debitPending, account.CreditsPending, creditPosted))
	available := creditPosted - debitPosted - debitPending
	reserved := debitPending
	balance := AccountBalance{Available: available, Reserved: reserved}
	return LedgerAccount{Id: account.ID.String(), Balance: balance}
}
