package ledger

import (
	"fmt"

	"encore.app/ledger/repository"
	"encore.app/ledger/service"
	"encore.app/ledger/workflow"
	"encore.dev"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	envName   = encore.Meta().Environment.Name
	taskQueue = envName + "-taskQueue"
)

//encore:service
type Service struct {
	ledgerService service.LedgerService
	client        client.Client
	worker        worker.Worker
}

func initService() (*Service, error) {

	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, taskQueue, worker.Options{})

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}

	tbLedger := repository.TBLedgerRepository{}
	tbLedger.Init()
	service := service.LedgerService{LedgerRepo: &tbLedger}

	// Temporal Workflow
	w.RegisterWorkflow(workflow.Greeting)
	w.RegisterActivity(workflow.ComposeGreeting)

	w.RegisterWorkflow(workflow.ExpireAuthorization)
	activities := &workflow.Activities{LedgerService: service}
	w.RegisterActivity(activities)

	return &Service{ledgerService: service, client: c, worker: w}, nil
}
