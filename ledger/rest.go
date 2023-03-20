package ledger

import (
	"context"

	"encore.app/ledger/service"
	"encore.app/ledger/workflow"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"go.temporal.io/sdk/client"
)

// Get Account Balance
//
// encore:api public method=GET path=/ledger/:id/balance
func (s *Service) GetAccountBalance(ctx context.Context, id string) (service.LedgerAccount, error) {
	options := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}

	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.Greeting, "asdasd")
	if err != nil {
		rlog.Error("error workflow")
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

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
	authorizationId, err := s.ledgerService.AuthorizePayment(id, payload.Amount)
	if err != nil {
		return AuthorizePaymentResponse{""}, err
	}

	s.workflowService.RunWF(ctx, id, payload.Amount, authorizationId, taskQueue, workflow.ExpireAuthorization)
	s.workflowService.SearchExpirationWF(ctx, id, payload.Amount)
	return AuthorizePaymentResponse{authorizationId}, nil
}

type PresentmentPayload struct {
	Amount int
}

// Presenment
//
// encore:api public method=POST path=/ledger/:id/presentment
func (s *Service) Presentment(ctx context.Context, id string, payload *PresentmentPayload) (AuthorizePaymentResponse, error) {

	wfId := s.workflowService.SearchExpirationWF(ctx, id, payload.Amount)

	// Factory would be nice here
	// Presentment Without Authorisation
	if len(wfId) <= 0 {
		id, _ := s.ledgerService.CreatePayment(id, payload.Amount)
		return AuthorizePaymentResponse{id}, nil
	}
	// Presentment With Authorisation
	s.client.SignalWorkflow(ctx, wfId, "", "cancel", "")
	authorizationId, _ := s.ledgerService.ReleasePayment(wfId)

	rlog.Info(authorizationId)
	return AuthorizePaymentResponse{authorizationId}, nil
}
