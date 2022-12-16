package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/pkg/errors"
)

const (
	TypeEmailWelcome = "email:welcome"
)

type EmailWelcomePayload struct {
	Username string
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewEmailWelcomeTask(username string) error {
	payload, err := json.Marshal(EmailWelcomePayload{Username: username})
	if err != nil {
		return errors.Wrapf(err, "[tasks] json marshal error, name: %s", TypeEmailWelcome)
	}
	task := asynq.NewTask(TypeEmailWelcome, payload)
	info, err := GetClient().Enqueue(task)
	if err != nil {
		return errors.Wrapf(err, "[tasks] Enqueue task error, name: %s", TypeEmailWelcome)
	}

	log.Infof("[tasks] welcome task", "info", info)

	return nil
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

func HandleEmailWelcomeTask(ctx context.Context, t *asynq.Task) error {
	var p EmailWelcomePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Infof("Sending Email to User: username=%s", p.Username)
	// Email delivery code ...
	return nil
}
