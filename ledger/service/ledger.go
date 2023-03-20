package service

import (
	"encore.app/ledger/repository"
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

func (s *LedgerService) AuthorizePayment(customerId string, amount int) (string, error) {
	transferId, _ := s.LedgerRepo.CreatePendingTransfer(customerId, amount)
	return transferId, nil
}

func (s *LedgerService) ReleasePayment(customerId string, amount int) (string, error) {
	// Does matching using workflow
	// Branch here, like presentment logic
	id := "1679286557552265000"
	transferId, _ := s.LedgerRepo.PostPendingTransfer(id)
	return transferId, nil
}

func (s *LedgerService) VoidPendingPayment(id string) (string, error) {
	// Does matching using workflow
	// id := "1679286557552265000"
	transferId, _ := s.LedgerRepo.VoidPendingTransfer(id)
	return transferId, nil
}

// Assume Debit is asset and Credit is liability
// Account 0 is the bank, others are customers
// So balance is technically credit
// This will cause integer overflow for bank typed account
func buildLedgerAccount(account *tb_types.Account) LedgerAccount {
	creditPosted := account.CreditsPosted
	debitPosted, debitPending := account.DebitsPosted, account.DebitsPending

	available := creditPosted - debitPosted
	reserved := debitPending
	balance := AccountBalance{Available: available, Reserved: reserved}
	return LedgerAccount{Id: account.ID.String(), Balance: balance}
}
