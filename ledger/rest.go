package ledger

import (
	"context"
	"fmt"

	"encore.app/ledger/service"
	"encore.app/ledger/workflow"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// Get Account Balance
//
// encore:api public method=GET path=/ledger/:id/balance
func (s *Service) GetAccountBalance(ctx context.Context, id string) (service.LedgerAccount, error) {
	searchAttributes := map[string]interface{}{
		"TransactionPendingAmount": "100",
		"TransactionUserId":        2,
	}

	options := client.StartWorkflowOptions{
		TaskQueue:        taskQueue,
		SearchAttributes: searchAttributes,
	}

	r, _ := s.client.ListWorkflow(ctx, &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: "default",
		// Query:     "TestSearch='asdasd' order by StartTime desc",
		Query: "ExecutionStatus = 'Completed' and TransactionPendingAmount='100' and TransactionUserId=2",
	})

	for k, v := range r.GetExecutions() {
		rlog.Warn(string(k))
		rlog.Warn(fmt.Sprint(v.Execution.WorkflowId, "  ", v.Execution.RunId, v.StartTime))
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
	authorizationId, _ := s.ledgerService.AuthorizePayment(id, payload.Amount)
	// if err != nil {
	// 	return nil, errs.B().Cause(err).Err()
	// }
	rlog.Info(authorizationId)

	options := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}
	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.ExpireAuthorization, authorizationId)
	if err != nil {
		rlog.Error("error workflow")
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())

	return AuthorizePaymentResponse{authorizationId}, nil
}

type PresentmentPayload struct {
	Amount int
}

// Presenment
//
//encore:api public method=POST path=/ledger/:id/presentment
func (s *Service) Presentment(ctx context.Context, id string, payload *PresentmentPayload) (AuthorizePaymentResponse, error) {
	authorizationId, _ := s.ledgerService.ReleasePayment(id, payload.Amount)

	rlog.Info(authorizationId)
	return AuthorizePaymentResponse{authorizationId}, nil
}
