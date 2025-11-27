package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	genChan := generator(ctx, nums)
	transChan, transErrChan := transform(ctx, genChan)
	saveChan, saveErrChan := save(ctx, transChan)
	for {
		select {
		case err, ok := <-transErrChan:
			if !ok {
				transErrChan = nil
				continue
			}
			fmt.Printf("transform error: %v\n", err)
			cancel()
		case err, ok := <-saveErrChan:
			fmt.Printf("save error: %v\n", err)
			if !ok {
				saveErrChan = nil
				continue
			}
			cancel()
		case <-saveChan:
			fmt.Printf("finished processing\n")
			return
		}
	}
}

func generator(ctx context.Context, nums []int) <-chan int {
	outChan := make(chan int)
	go func() {
		defer close(outChan)
		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			case outChan <- num:
			}
		}
	}()
	return outChan
}

func transform(ctx context.Context, inChan <-chan int) (<-chan int, <-chan error) {
	outChan := make(chan int)
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		defer close(outChan)
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			case num, ok := <-inChan:
				if !ok {
					return
				}
				// if num == 6 {
				// 	errChan <- fmt.Errorf("number %d is invalid", num)
				// 	return
				// }
				outChan <- num * 2
			}
		}
	}()
	return outChan, errChan

}

func save(ctx context.Context, inChan <-chan int) (<-chan struct{}, <-chan error) {
	doneChan := make(chan struct{})
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		defer close(doneChan)
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			case num, ok := <-inChan:
				if !ok {
					return
				}
				fmt.Printf("saved %d\n", num)
			}
		}
	}()
	return doneChan, errChan
}
