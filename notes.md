# Notes

## Context
Many functions in Go use the context package to gather additional information about the environment they’re being executed in, and will typically provide that context to the functions they also call. When using contexts, it’s important to know that the values stored in a specific context.Context are immutable, meaning they can’t be changed. The context.WithValue function didn’t modify the context you provided. Instead, it wrapped your parent context inside another one with the new value.

A good rule of thumb is that any data required for a function to run should be passed as parameters. Sometimes, for example, it can be useful to keep values such as usernames in context values for use when logging information for later. However, if the username is used to determine if a function should display some specific information, you’d want to include it as a function parameter even if it’s already available from the context. This way when you, or someone else, looks at the function in the future, it’s easier to see which data is actually being used.

### Example
If a web server function is handling an HTTP request for a specific client, the function may only need to know which URL the client is requesting to serve the response. The function might only take that URL as a parameter. However, things can always happen when serving a response, such as the client disconnecting before receiving the response. If the function serving the response doesn’t know the client disconnected, the server software may end up spending more computing time than it needs calculating a response that will never be used.

In this case, being aware of the context of the request, such as the client’s connection status, allows the server to stop processing the request once the client disconnects. This saves valuable compute resources on a busy server and frees them up to handle another client’s request. This type of information can also be helpful in other contexts where functions take time to execute, such as making database calls. To enable ubiquitous access to this type of information, Go has included a context package in its standard library.

### Functions
context.TODO() - empty/starting context

context.Background() - empty/starting context but it's desinged to be used where you intend to start a known context

context.WithValue() - WithValue returns a copy of parent in which the value associated with key is val. Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.

context.WithCancel() - WithCancel returns a copy of parent with a new Done channel. The returned context's Done channel is closed when the returned cancel function is called or when the parent context's Done channel is closed, whichever happens first.

context.WithDeadline() - WithDeadline returns a copy of the parent context with the deadline adjusted to be no later than d. If the parent's deadline is already earlier than d, WithDeadline(parent, d) is semantically equivalent to parent. The returned [Context.Done] channel is closed when the deadline expires, when the returned cancel function is called, or when the parent context's Done channel is closed, whichever happens first.

## Select Statements
A "select" statement chooses which of a set of possible send or receive operations will proceed. It looks similar to a "switch" statement but with the cases all referring to communication operations. 

```
SelectStmt = "select" "{" { CommClause } "}" .
CommClause = CommCase ":" StatementList .
CommCase   = "case" ( SendStmt | RecvStmt ) | "default" .
RecvStmt   = [ ExpressionList "=" | IdentifierList ":=" ] RecvExpr .
RecvExpr   = Expression .
```

## Channels
Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values into another goroutine.

```
package main

import "fmt"

func main() {
    // Create a new channel with make(chan val-type). Channels are typed by the values they convey.
    messages := make(chan string)

    // Send a value into a channel using the channel <- syntax. Here we send "ping" to the messages channel we made above, from a new goroutine.
    go func() { messages <- "ping" }()

    // The <-channel syntax receives a value from the channel. Here we’ll receive the "ping" message we sent above and print it out.
    msg := <-messages
    fmt.Println(msg)
}
```

### Channel Buffering
By default channels are unbuffered, meaning that they will only accept sends (chan <-) if there is a corresponding receive (<- chan) ready to receive the sent value. Buffered channels accept a limited number of values without a corresponding receiver for those values.

```
package main	

import "fmt"

func main() {

    // Here we make a channel of strings buffering up to 2 values.
    messages := make(chan string, 2)

    // Because this channel is buffered, we can send these values into the channel without a corresponding concurrent receive.
    messages <- "buffered"
    messages <- "channel"

    // Later we can receive these two values as usual.
    fmt.Println(<-messages)
    fmt.Println(<-messages)
}
```

### Channel Synchronization
We can use channels to synchronize execution across goroutines. Here’s an example of using a blocking receive to wait for a goroutine to finish. When waiting for multiple goroutines to finish, you may prefer to use a WaitGroup.

```
package main
import (
    "fmt"
    "time"
)

// This is the function we’ll run in a goroutine. The done channel will be used to notify another goroutine that this function’s work is done.
func worker(done chan bool) {
    fmt.Print("working...")
    time.Sleep(time.Second)
    fmt.Println("done")

    // Send a value to notify that we’re done.
    done <- true
}

func main() {
    // Start a worker goroutine, giving it the channel to notify on.
    done := make(chan bool, 1)
    go worker(done)

    Block until we receive a notification from the worker on the channel.
    <-done
}
```

### Channel Directions
When using channels as function parameters, you can specify if a channel is meant to only send or receive values. This specificity increases the type-safety of the program.

```
package main
import "fmt"

// This ping function only accepts a channel for sending values. It would be a compile-time error to try to receive on this channel.
func ping(pings chan<- string, msg string) {
    pings <- msg
}

// The pong function accepts one channel for receives (pings) and a second for sends (pongs).
func pong(pings <-chan string, pongs chan<- string) {
    msg := <-pings
    pongs <- msg
}

func main() {
    pings := make(chan string, 1)
    pongs := make(chan string, 1)
    ping(pings, "passed message")
    pong(pings, pongs)
    fmt.Println(<-pongs)
}
```