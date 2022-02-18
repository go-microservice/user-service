package main

import (
	"log"
	"time"

	"github.com/go-microservice/user-service/internal/task"

	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func init() {
	task.Init()
}

func main() {
	// run worker server
	go func() {
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: redisAddr},
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: 10,
				// Optionally specify multiple queues with different priority.
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
				// See the godoc for other configuration options
			},
		)

		// mux maps a type to a handler
		mux := asynq.NewServeMux()
		// register handlers...
		mux.HandleFunc(task.TypeEmailWelcome, task.HandleEmailWelcomeTask)
		//mux.Handle(task.TypeImageResize, task.NewImageProcessor())

		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// run schedule server
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		&asynq.SchedulerOpts{Location: time.Local},
	)

	t, _ := task.NewEmailWelcomeTask(555)
	if _, err := scheduler.Register("@every 5s", t); err != nil {
		log.Fatal(err)
	}

	// Run blocks and waits for os signal to terminate the program.
	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
