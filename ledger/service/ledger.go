package service

import (
	"context"

	"encore.app/ledger/repository"
	"encore.app/ledger/workflow"
	tb_types "github.com/tigerbeetledb/tigerbeetle-go/pkg/types"
)

type LedgerService struct {
	LedgerRepo      repository.LedgerRepository
	WorkflowService WorkflowService
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

func (s *LedgerService) AuthorizePayment(ctx context.Context, customerId string, amount int) (string, error) {
	transferId, err := s.LedgerRepo.CreatePendingTransfer(customerId, amount)

	s.WorkflowService.RunWF(ctx, customerId, amount, transferId, workflow.ExpireAuthorization)

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

func (s *LedgerService) Presentment(ctx context.Context, id string, amount int) (string, error) {
	wfId := s.WorkflowService.SearchExpirationWF(ctx, id, amount)

	if len(wfId) <= 0 {
		id, err := s.CreatePayment(id, amount)

		if err != nil {
			return "", err
		}

		return id, nil
	}

	// Presentment With Authorisation
	s.WorkflowService.CancelExpirationWF(ctx, wfId)
	authorizationId, _ := s.ReleasePayment(wfId)

	return authorizationId, nil
}

// Assume Debit is asset and Credit is liability
// Account 1 is the bank, others are customers
// So balance is customer POV is technically credit
// TODO: Handle overflow
func buildLedgerAccount(account *tb_types.Account) LedgerAccount {
	creditPosted := account.CreditsPosted
	debitPosted, debitPending := account.DebitsPosted, account.DebitsPending

	available := creditPosted - debitPosted - debitPending
	reserved := debitPending
	balance := AccountBalance{Available: available, Reserved: reserved}
	return LedgerAccount{Id: account.ID.String(), Balance: balance}
}
