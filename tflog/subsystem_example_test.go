package tflog

import "os"

func ExampleNewSubsystem() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// non-example-setup code begins here

	// register a new subsystem before using it
	subCtx := NewSubsystem(exampleCtx, "my-subsystem")

	// messages logged to the subsystem will carry the subsystem name with
	// them
	SubsystemTrace(subCtx, "my-subsystem", "hello, world", "foo", 123)

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"provider.my-subsystem","foo":123}
}

func ExampleNewSubsystem_withLevel() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// simulate the user setting a logging level for the subsystem
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
	// {"@level":"trace","@message":"hello, world","@module":"provider","foo":123}
	// {"@level":"warn","@message":"hello, world","@module":"provider.my-subsystem","foo":123}
}

func ExampleSubsystemWith() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// register a new subsystem before using it
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here

	// associate a key and value with all lines logged by the sub-logger
	derivedCtx := SubsystemWith(exampleCtx, "my-subsystem", "foo", 123)

	// all messages logged with derivedCtx will now have foo=123
	// automatically included
	SubsystemTrace(derivedCtx, "my-subsystem", "example log message")

	// Output:
	// {"@level":"trace","@message":"example log message","@module":"provider.my-subsystem","foo":123}
}

func ExampleSubsystemTrace() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// register a new subsystem before using it
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemTrace(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"trace","@message":"hello, world","@module":"provider.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemDebug() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// register a new subsystem before using it
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemDebug(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"debug","@message":"hello, world","@module":"provider.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemInfo() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// register a new subsystem before using it
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemInfo(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"info","@message":"hello, world","@module":"provider.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemWarn() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// register a new subsystem before using it
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemWarn(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"warn","@message":"hello, world","@module":"provider.my-subsystem","colors":["red","blue","green"],"foo":123}
}

func ExampleSubsystemError() {
	// virtually no plugin developers will need to worry about
	// instantiating loggers, as the libraries they're using will take care
	// of that, but we're not using those libraries in these examples. So
	// we need to do the injection ourselves. Plugin developers will
	// basically never need to do this, so the next line can safely be
	// considered setup for the example and ignored. Instead, use the
	// context passed in by the framework or library you're using.
	exampleCtx := getExampleContext()

	// register a new subsystem before using it
	exampleCtx = NewSubsystem(exampleCtx, "my-subsystem")

	// non-example-setup code begins here
	SubsystemError(exampleCtx, "my-subsystem", "hello, world", "foo", 123, "colors", []string{"red", "blue", "green"})

	// Output:
	// {"@level":"error","@message":"hello, world","@module":"provider.my-subsystem","colors":["red","blue","green"],"foo":123}
}
