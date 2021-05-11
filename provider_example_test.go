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
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
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
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Trace(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"trace","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleDebug() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Debug(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"debug","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleInfo() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Info(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"info","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleWarn() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Warn(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"warn","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}

func ExampleError() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Error(exampleCtx, "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"error","@message":"hello, world","colors":["red","blue","green"],"foo":123}
}
