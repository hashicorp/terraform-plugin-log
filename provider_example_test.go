package tflog

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
)

func getExampleContext() context.Context {
	return New(context.Background(), withOutput(os.Stdout),
		WithLevel(hclog.Trace), WithoutLocation(), withoutTimestamp())
}

func ExampleWith() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	derivedCtx := With(exampleCtx, "foo", 123)

	// all messages logged with derivedCtx will now have foo=123
	// automatically included
	Trace(derivedCtx, "example log message")

	// Output:
	// {"@level":"trace","@message":"example log message","foo":123}
}

func ExampleTrace() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Trace(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"trace","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleDebug() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Debug(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"debug","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleInfo() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Info(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"info","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleWarn() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Warn(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"warn","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleError() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Error(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"error","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}
