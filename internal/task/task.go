package task

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeEmailWelcome = "email:deliver"
	TypeImageResize  = "image:resize"
)

const redisAddr = "127.0.0.1:6379"

func Init() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Enqueue task to be processed immediately.
	// Use (*Client).Enqueue method.
	// ------------------------------------------------------
	task, err := NewEmailWelcomeTask(111)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Schedule task to be processed in the future.
	// Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------
	task, err = NewEmailWelcomeTask(222)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.ProcessIn(10*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Set other options to tune task processing behavior.
	// Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------
	task, err = NewEmailWelcomeTask(333)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// scheduler task, like crontab
	// ----------------------------------------------------------------------------

}
