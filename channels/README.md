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
