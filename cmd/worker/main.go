package main

import (
	"context"
	"fmt"
	"github.com/tgmendes/soundfuse/worker"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type WorkerReq struct {
	JobID    string
	Duration time.Duration
}

var (
	reqChan chan WorkerReq
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	// serverErrors := make(chan error, 1)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wPool := worker.Pool{
		TasksChan:  make(chan *worker.Task, 5),
		NumWorkers: 1,
	}

	wPool.Run(ctx)
	for i := 0; i < 20; i++ {
		i := i
		taskID := uuid.NewString()

		tf := func() error {
			res, err := DoSomething(taskID, i)
			fmt.Printf("Got result from processing: %s\n", res)
			return err
		}
		wPool.AddTask(&worker.Task{ReqID: taskID, F: tf})
	}
	fmt.Println("added tasks")
	time.Sleep(5 * time.Second)

	tf := func() error {
		fmt.Println("impostor")
		return nil
	}
	wPool.AddTask(&worker.Task{ReqID: "rando", F: tf})
	time.Sleep(5 * time.Second)
	cancel()

	time.Sleep(2 * time.Second)
	// reqChan = make(chan WorkerReq)
	// r := chi.NewRouter()
	// r.Get("/", workHandler)
	//
	// for w := 1; w <= 5; w++ {
	// 	go func(ctx context.Context, w int) {
	// 		worker(ctx, w, reqChan)
	// 	}(ctx, w)
	//
	// }
	// time.Sleep(2 * time.Second)
	// cancel()
	// time.Sleep(2 * time.Second)

	// // Start the service listening for api requests.
	// go func() {
	// 	serverErrors <- http.ListenAndServe(":8082", r)
	// }()
	//
	// // =========================================================================
	// // Shutdown
	//
	// // Blocking main and waiting for shutdown.
	// select {
	// case err := <-serverErrors:
	// 	log.Fatalf("server encountered error: %s\n", err)
	//
	// case sig := <-shutdown:
	// 	log.Printf("initialising shutdown: %s\n", sig)
	// 	// Give outstanding requests a deadline for completion.
	// 	cancel()
	// 	close(reqChan)
	// 	log.Printf("shutdown complete: %s\n", sig)
	// }
}

func DoSomething(reqID string, idx int) (string, error) {
	dur := rand.Intn(5-1) + 1
	fmt.Printf("Processing request #%d with ID %s. Sleeping %d ms\n", idx, reqID, dur*100)
	time.Sleep(time.Duration(dur) * time.Millisecond * 100)
	fmt.Printf("finished processing ID %s\n", reqID)
	if idx == 7 || idx == 12 || idx == 15 {
		return "", fmt.Errorf("failed on idx %d\n", idx)
	}
	return reqID, nil
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewString()
	dur := rand.Intn(5-1) + 1
	reqChan <- WorkerReq{JobID: id, Duration: time.Second * 10 * time.Duration(dur)}

	w.Write([]byte(fmt.Sprintf("worker: %s", id)))
}
