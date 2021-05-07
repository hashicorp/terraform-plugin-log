package tfsdklog

import "os"

func ExampleNewSubsystem() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()

	// non-example-setup code begins here
	subCtx := NewSubsystem(exampleCtx, "my-subsystem")

	// messages logged to the subsystem will carry the subsystem name with
	// them
	SubsystemTrace(subCtx, "my-subsystem", "hello, world", "foo", 123)

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"sdk.my-subsystem","foo":123}
}

func ExampleNewSubsystem_withLevel() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	os.Setenv("EXAMPLE_SUBSYSTEM_LEVEL", "WARN")

	// non-example-setup code begins here

	// create a context with a logger for a new "my-subsystem" subsystem,
	// using the WARN level from the "EXAMPLE_SUBSYSTEM_LEVEL" environment
	// variable
	subCtx := NewSubsystem(exampleCtx, "my-subsystem", WithLevelFromEnv("EXAMPLE_SUBSYSTEM_LEVEL"))

	// this won't actually get output, it's not at WARN or higher
	SubsystemTrace(subCtx, "my-subsystem", "hello, world", "foo", 123)

	// the parent logger will still output at its configured TRACE level,
	// though
	Trace(subCtx, "hello, world", "foo", 123)

	// and the subsystem logger will output at the WARN level
	SubsystemWarn(subCtx, "my-subsystem", "hello, world", "foo", 123)

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"sdk","foo":123}
	// {"@level":"warn","@message":"hello, world","@module":"sdk.my-subsystem","foo":123}
}

func ExampleSubsystemWith() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	derivedCtx := SubsystemWith(exampleCtx, "my-subsystem", "foo", 123)

	// all messages logged with derivedCtx will now have foo=123
	// automatically included
	SubsystemTrace(derivedCtx, "my-subsystem", "example log message")

	// Output:
	// {"@level":"trace","@message":"example log message","@module":"sdk.my-subsystem","foo":123}
}

func ExampleSubsystemTrace() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemTrace(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"sdk.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemDebug() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemDebug(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"debug","@message":"hello, world","@module":"sdk.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemInfo() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemInfo(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"info","@message":"hello, world","@module":"sdk.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemWarn() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemWarn(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"warn","@message":"hello, world","@module":"sdk.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemError() {
	// context for example only; plugins should never need to do this
	exampleCtx := getExampleContext()
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemError(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"error","@message":"hello, world","@module":"sdk.my-subsystem","colors":["red","blue","green"],"foo":123}
}
