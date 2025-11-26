package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(signalCtx)
	defer cancel()

	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	generatorChan := generator(ctx, nums)
	transformChan, transformErrorChan := transform(ctx, generatorChan)
	saveDoneChan, saveErrorChan := save(ctx, transformChan)

	mergedErrorChans := mergeErrorChannels(ctx, transformErrorChan, saveErrorChan)

	for {
		select {
		case err, ok := <-mergedErrorChans:
			if !ok {
				mergedErrorChans = nil
				continue
			}
			fmt.Printf("pipeline error: %v\n", err)
			cancel()
			return
		case <-saveDoneChan:
			fmt.Printf("successfully finished processing\n")
			return
		}
	}
}

func mergeErrorChannels(ctx context.Context, errChansToMerge ...<-chan error) <-chan error {
	merged := make(chan error)
	var wg sync.WaitGroup
	for _, errChan := range errChansToMerge {
		wg.Add(1)
		go func(ch <-chan error) {
			defer wg.Done()
			for err := range ch {
				select {
				case <-ctx.Done():
					return
				case merged <- err:
				}
			}
		}(errChan)
	}
	go func() {
		wg.Wait()
		close(merged)
	}()
	return merged
}

func generator(ctx context.Context, nums []int) <-chan int {
	outChan := make(chan int)
	go func() {
		defer close(outChan)
		for _, num := range nums {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
				outChan <- num
			}
		}
	}()
	return outChan
}

func transform(ctx context.Context, inChan <-chan int) (<-chan int, <-chan error) {
	errChan := make(chan error)
	outChan := make(chan int)

	go func() {
		defer close(outChan)
		defer close(errChan)
		for num := range inChan {
			select {
			case <-ctx.Done():
				return
			default:
				if num == 6 {
					errChan <- fmt.Errorf("transform error: number %d is invalid", num)
					return
				}
				time.Sleep(100 * time.Millisecond)
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
		defer close(doneChan)
		defer close(errChan)
		for num := range inChan {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("saved: %d\n", num)
			}
		}
	}()
	return doneChan, errChan
}
