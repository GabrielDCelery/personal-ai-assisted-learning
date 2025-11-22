# Learning channels with AI

## Lesson 1: Basic Channel Creation and Communication

Concept: Channels are typed conduits for communication between goroutines. They follow the principle "Don't communicate by sharing memory; share memory by communicating."

Your task:

1. Create an unbuffered channel of type string
2. Launch a goroutine that sends the message "hello from goroutine" to the channel
3. In main(), receive from the channel and print the result

Hints:

- Create a channel: ch := make(chan string)
- Send to channel: ch <- value
- Receive from channel: value := <-ch
- Launch goroutine: go func() { ... }()

Expected output:
hello from goroutine

Key learning: Unbuffered channels block on send until someone receives, and block on receive until someone sends. This provides synchronization.

## Lesson 2: Buffered Channels

Concept: Buffered channels have capacity. Sends only block when the buffer is full, receives only block when empty.

Your task:

1. Create a buffered channel of type int with capacity 3
2. Send three values (1, 2, 3) to the channel without using a goroutine
3. Receive and print all three values

Hints:

- Buffered channel: ch := make(chan int, 3)
- You can send up to 3 values before blocking

Expected output:
1
2
3

Key learning: With unbuffered channels, you couldn't do this without a goroutine (deadlock). Buffered channels decouple send/receive timing up to the buffer size.

## Lesson 3: Directional Channels (Channel Types)

Concept: You can restrict a channel to send-only or receive-only. This enforces correct usage at compile time and makes function signatures self-documenting.

Your task:

1. Create a function producer(ch chan<- int) that sends numbers 1-5 to the channel, then closes it
2. Create a function consumer(ch <-chan int) that receives all values and prints them
3. In main(), create a bidirectional channel, pass it to both functions (producer in a goroutine), and let consumer run

Hints:

- chan<- int = send-only (arrow points INTO channel)
- <-chan int = receive-only (arrow points OUT OF channel)
- A regular chan int can be passed where directional channels are expected (Go converts automatically)

Expected output:
1
2
3
4
5

Key learning: Directional channels prevent bugs. If consumer accidentally tries to close or send to the channel, the compiler catches it. This is especially valuable in larger codebases.

## Lesson 4: The select Statement

Concept: select lets you wait on multiple channel operations simultaneously. It's like a switch for channels - whichever case is ready first executes.

Your task:

1. Create two channels: ch1 and ch2 (both chan string)
2. Launch two goroutines:
   - First sends "from channel 1" after 100ms delay
   - Second sends "from channel 2" after 200ms delay

3. Use a select inside a loop to receive from whichever channel is ready first
4. Print both messages as they arrive, then exit

Hints:

- time.Sleep(100 \* time.Millisecond) for delays
- Basic select structure:

```go
  select {
  case msg := <-ch1:
  // handle ch1
  case msg := <-ch2:
  // handle ch2
  }
```

- You need to receive exactly 2 messages total

Expected output:
from channel 1
from channel 2

Key learning: select is fundamental for handling multiple concurrent operations, timeouts, and cancellation patterns.

### Why is a closed channel "always ready"?

This is a design decision in Go. A closed channel can always be read from - it returns the zero value instantly. This allows consumers to drain remaining buffered values and detect closure.

ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

fmt.Println(<-ch) // 1 (buffered value)
fmt.Println(<-ch) // 2 (buffered value)
fmt.Println(<-ch) // 0 (zero value, channel closed)
fmt.Println(<-ch) // 0 (zero value, forever)

Closing doesn't "lock" the channel - it signals "no more sends." Reads always succeed.

2. Why doesn't default get picked randomly?

default is not part of the random selection. The rules are:

1. If one or more channel cases are ready → pick randomly among only those ready cases
2. If zero channel cases are ready → run default
3. If zero channel cases are ready AND no default → block

default is the fallback for "nothing ready," not an equal participant.

3. How does ch2 ever get selected if ch1 is always ready?

Both closed channels are always ready simultaneously. So select randomly picks between ch1 and ch2 each iteration. Both get hit frequently, but default never runs because at least one case is always ready.

Iteration 1: ch1 ready, ch2 ready → random pick → ch1
Iteration 2: ch1 ready, ch2 ready → random pick → ch2
Iteration 3: ch1 ready, ch2 ready → random pick → ch1
... (default never runs)
