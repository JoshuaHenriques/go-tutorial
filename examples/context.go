package main

import (
	"context"
	"fmt"
	"time"
)

func doSomething(ctx context.Context) {
	// ctx, cancelCtx := context.WithCancel(ctx)

	// deadline := time.Now().Add(1500 * time.Millisecond)
	// ctx, cancelCtx := context.WithDeadline(ctx, deadline)

	ctx, cancelCtx := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancelCtx()

	printCh := make(chan int)
	go doAnother(ctx, printCh)

	// for num := 1; num <= 3; num++ {
	// 	printCh <- num
	// }
	for num := 1; num <= 3; num++ {
		select {
		case printCh <- num:
			time.Sleep(1 * time.Second)
		case <-ctx.Done():
			break
		}
	}

	cancelCtx()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("doSomething: finished\n")

	// fmt.Println("Doing something")
	// fmt.Printf("doSomething: myKey's value is %s\n", ctx.Value("myKey"))

	// anotherCtx := context.WithValue(ctx, "myKey", "anotherValue")
	// doAnother(anotherCtx)

	// fmt.Printf("doSomething: myKey's value is %s\n", ctx.Value("myKey"))
}

func doAnother(ctx context.Context, printCh <-chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("doAnother err: %s\n", err)
			}
			fmt.Printf("doAnother: finished\n")
			return
		case num := <-printCh:
			fmt.Printf("doAnother: %d\n", num)
		}
	}
	// fmt.Printf("doAnother: myKey's value is %s\n", ctx.Value("myKey"))
}

func main() {
	// context.TODO() - empty/starting context
	// context.Background() - empty/starting context but it's desinged to be used where
	// you intend to start a known context
	ctx := context.Background()

	ctx = context.WithValue(ctx, "myKey", "myValue")
	doSomething(ctx)
}
