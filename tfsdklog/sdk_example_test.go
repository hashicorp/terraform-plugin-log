package tfsdklog

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

func getExampleContext() context.Context {
	return NewRootSDKLogger(context.Background(), logging.WithOutput(os.Stdout),
		WithLevel(hclog.Trace), WithoutLocation(), logging.WithoutTimestamp())
}

func ExampleSetField() {
	// this function calls New with the options it needs to be reliably
	// tested. Framework and SDK developers should call New, inject the
	// resulting context in their framework, and then pass it around. This
	// exampleCtx is a stand-in for a context you have injected a logger
	// into and passed to the area of the codebase you need it.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	derivedCtx := SetField(exampleCtx, "foo", 123)

	// all messages logged with derivedCtx will now have foo=123
	// automatically included
	Trace(derivedCtx, "example log message")

	// Output:
	// {"@level":"trace","@message":"example log message","@module":"sdk","foo":123}
}

func ExampleTrace() {
	// this function calls New with the options it needs to be reliably
	// tested. Framework and SDK developers should call New, inject the
	// resulting context in their framework, and then pass it around. This
	// exampleCtx is a stand-in for a context you have injected a logger
	// into and passed to the area of the codebase you need it.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Trace(exampleCtx, "hello, world", map[string]interface{}{
		"foo":    123,
		"colors": []string{"red", "blue", "green"},
	})

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"sdk","colors":["red","blue","green"],"foo":123}
}

func ExampleDebug() {
	// this function calls New with the options it needs to be reliably
	// tested. Framework and SDK developers should call New, inject the
	// resulting context in their framework, and then pass it around. This
	// exampleCtx is a stand-in for a context you have injected a logger
	// into and passed to the area of the codebase you need it.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Debug(exampleCtx, "hello, world", map[string]interface{}{
		"foo":    123,
		"colors": []string{"red", "blue", "green"},
	})

	// Output:
	// {"@level":"debug","@message":"hello, world","@module":"sdk","colors":["red","blue","green"],"foo":123}
}

func ExampleInfo() {
	// this function calls New with the options it needs to be reliably
	// tested. Framework and SDK developers should call New, inject the
	// resulting context in their framework, and then pass it around. This
	// exampleCtx is a stand-in for a context you have injected a logger
	// into and passed to the area of the codebase you need it.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Info(exampleCtx, "hello, world", map[string]interface{}{
		"foo":    123,
		"colors": []string{"red", "blue", "green"},
	})

	// Output:
	// {"@level":"info","@message":"hello, world","@module":"sdk","colors":["red","blue","green"],"foo":123}
}

func ExampleWarn() {
	// this function calls New with the options it needs to be reliably
	// tested. Framework and SDK developers should call New, inject the
	// resulting context in their framework, and then pass it around. This
	// exampleCtx is a stand-in for a context you have injected a logger
	// into and passed to the area of the codebase you need it.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Warn(exampleCtx, "hello, world", map[string]interface{}{
		"foo":    123,
		"colors": []string{"red", "blue", "green"},
	})

	// Output:
	// {"@level":"warn","@message":"hello, world","@module":"sdk","colors":["red","blue","green"],"foo":123}
}

func ExampleError() {
	// this function calls New with the options it needs to be reliably
	// tested. Framework and SDK developers should call New, inject the
	// resulting context in their framework, and then pass it around. This
	// exampleCtx is a stand-in for a context you have injected a logger
	// into and passed to the area of the codebase you need it.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	Error(exampleCtx, "hello, world", map[string]interface{}{
		"foo":    123,
		"colors": []string{"red", "blue", "green"},
	})

	// Output:
	// {"@level":"error","@message":"hello, world","@module":"sdk","colors":["red","blue","green"],"foo":123}
}
