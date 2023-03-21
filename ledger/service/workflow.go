package service

import (
	"context"
	"fmt"

	"encore.dev/rlog"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type WorkflowService struct {
	Client    client.Client
	TaskQueue string
}

func (s *WorkflowService) SearchExpirationWF(ctx context.Context, id string, amount int) string {
	query := fmt.Sprintf("ExecutionStatus = 'Running' and TransactionPendingAmount=%d and TransactionUserId='%s'", amount, id)
	r, _ := s.Client.ListWorkflow(ctx, &workflowservice.
		ListWorkflowExecutionsRequest{
		Namespace: "default",
		Query:     query,
	})

	wfArr := r.GetExecutions()
	wfArrLength := len(wfArr)
	if wfArrLength <= 0 {
		return ""
	}

	for _, v := range r.GetExecutions() {
		rlog.Warn(fmt.Sprint(v.Execution.WorkflowId, "  ", v.Execution.RunId, v.StartTime))
	}

	// Get Oldest Match
	wfId := wfArr[wfArrLength-1].Execution.GetWorkflowId()
	rlog.Info(fmt.Sprint("Found: ", wfId))
	return wfId
}

func (s *WorkflowService) RunWF(ctx context.Context, id string, amount int, authorizationId string, workflow interface{}) {
	searchAttributes := map[string]interface{}{
		"TransactionPendingAmount": amount,
		"TransactionUserId":        id,
	}

	options := client.StartWorkflowOptions{
		ID:               authorizationId,
		TaskQueue:        s.TaskQueue,
		SearchAttributes: searchAttributes,
	}

	we, err := s.Client.ExecuteWorkflow(ctx, options, workflow, authorizationId)
	if err != nil {
		rlog.Error("error workflow")
	}
	rlog.Info("started workflow", "id", we.GetID(), "run_id", we.GetRunID())
}

func (s *WorkflowService) CancelExpirationWF(ctx context.Context, wfId string) {
	s.Client.SignalWorkflow(ctx, wfId, "", "cancel", "")
}
