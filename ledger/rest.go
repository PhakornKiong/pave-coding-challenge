package ledger

import (
	"context"

	"encore.app/ledger/service"
	"encore.dev/beta/errs"
)

// Get Account Balance
//
// encore:api public method=GET path=/ledger/:id/balance
func (s *Service) GetAccountBalance(ctx context.Context, id string) (service.LedgerAccount, error) {
	acc, err := s.ledgerService.GetAccountBalance(id)
	return acc, err
}

type CreateAccountPayload struct {
	Id string
}

// Create Account
//
//encore:api public method=POST path=/ledger
func (s *Service) CreateAccount(ctx context.Context, payload *CreateAccountPayload) error {
	err := s.ledgerService.CreateAccount(payload.Id)
	if err != nil {
		return errs.B().Cause(err).Err()
	}
	return nil
}

type AddBalancePayload struct {
	Amount int
}

// Add Account Balance
//
//encore:api public method=POST path=/ledger/:id/addBalance
func (s *Service) AddBalance(ctx context.Context, id string, payload *AddBalancePayload) error {
	_, err := s.ledgerService.AddAccountBalance(id, payload.Amount)
	if err != nil {
		return errs.B().Cause(err).Err()
	}
	return nil
}

type AuthorizePaymentPayload struct {
	Amount int
}

type AuthorizePaymentResponse struct {
	Id string
}

// Authorize Payment
//
//encore:api public method=POST path=/ledger/:id/authorize
func (s *Service) AuthorizePayment(ctx context.Context, id string, payload *AuthorizePaymentPayload) (AuthorizePaymentResponse, error) {
	authorizationId, err := s.ledgerService.AuthorizePayment(ctx, id, payload.Amount)

	if err != nil {
		return AuthorizePaymentResponse{""}, err
	}

	return AuthorizePaymentResponse{authorizationId}, nil
}

type PresentmentPayload struct {
	Amount int
}

// Presenment
//
// encore:api public method=POST path=/ledger/:id/presentment
func (s *Service) Presentment(ctx context.Context, id string, payload *PresentmentPayload) (AuthorizePaymentResponse, error) {

	id, err := s.ledgerService.Presentment(ctx, id, payload.Amount)

	if err != nil {
		return AuthorizePaymentResponse{""}, err
	}

	return AuthorizePaymentResponse{id}, nil
}
