package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go Task(ctx, "worker 1")

	go Task(ctx, "worker 2")

	go Task(ctx, "worker 3")

	time.Sleep(time.Second * 3)
	cancel()

	time.Sleep(time.Second * 1)
}

func Task(ctx context.Context, taskName string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task is canceled..")
			return
		default:
			fmt.Printf("%s is running..\n", taskName)
			time.Sleep(time.Second * 1)
		}
	}

}
