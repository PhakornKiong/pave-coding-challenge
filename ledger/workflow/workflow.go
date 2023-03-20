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

func ExpireAuthorization(ctx workflow.Context, id string) (string, error) {
	var a *Activities
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 1000,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	workflow.Sleep(ctx, 5*time.Second)

	var result string
	err := workflow.ExecuteActivity(ctx, a.ExpireAuthorization, id).Get(ctx, &result)

	return result, err
}
