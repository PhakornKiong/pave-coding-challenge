package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	expirationTimeout = time.Second * 100
)

func ExpireAuthorization(ctx workflow.Context, id string) error {

	var a *Activities
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 1000,
	}
	cancelChan := workflow.GetSignalChannel(ctx, "cancel")

	ctx = workflow.WithActivityOptions(ctx, options)

	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)

	var result string
	timerFuture := workflow.NewTimer(childCtx, expirationTimeout)

	selector.AddReceive(cancelChan, func(c workflow.ReceiveChannel, more bool) {
		// If a cancel signal is received, cancel workflow
		workflow.GetLogger(ctx).Info("Cancel signal received")
		c.Receive(ctx, nil)
		cancelHandler()
	})

	selector.AddFuture(timerFuture, func(f workflow.Future) {
		workflow.ExecuteActivity(ctx, a.ExpireAuthorization, id).Get(ctx, &result)
	})

	selector.Select(ctx)

	workflow.GetLogger(ctx).Info("Expiration Workflow completed.")
	return nil
}
