package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func Greeting(ctx workflow.Context, name string) (string, error) {

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	err := workflow.ExecuteActivity(ctx, ComposeGreeting, name).Get(ctx, &result)

	return result, err
}

func ExpireAuthorization(ctx workflow.Context, id string) error {
	expirationTimeout := time.Second * 100
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
