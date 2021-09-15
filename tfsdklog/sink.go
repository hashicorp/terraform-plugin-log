package tfsdklog

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

const (
	// envLog is the environment variable that users can set to control the
	// least-verbose level of logs that will be output during testing. If
	// this environment variable is set, it will default to off. This is
	// just the default; specific loggers and sub-loggers can set a lower
	// or higher verbosity level without a problem right now. In theory,
	// they should not be able to.
	//
	// Valid values are TRACE, DEBUG, INFO, WARN, ERROR, and OFF. A special
	// pseudo-value, JSON, will set the value to TRACE and output the
	// results in their JSON format.
	envLog = "TF_LOG"

	// envLogFile is the environment variable that controls where log
	// output is written during tests. By default, logs will be written to
	// standard error. Setting this environment variable to another file
	// path will write logs there instead during tests.
	envLogFile = "TF_LOG_PATH"
)

var sink hclog.Logger

func init() {
	sink = newSink()
}

// ValidLevels are the string representations of levels that can be set for
// loggers.
var ValidLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}

func getSink(ctx context.Context) hclog.Logger {
	logger := ctx.Value(logging.SinkKey)
	if logger == nil {
		return nil
	}
	return logger.(hclog.Logger)
}

// RegisterSink sets up a logging sink, for use with test frameworks and other
// cases where plugin logs don't get routed through Terraform. This applies the
// same filtering and file output behaviors that Terraform does.
//
// RegisterSink should only ever be called by test frameworks, providers should
// never call it.
//
// RegisterSink must be called prior to any loggers being setup or instantiated.
func RegisterSink(ctx context.Context) context.Context {
	return context.WithValue(ctx, logging.SinkKey, sink)
}

func newSink() hclog.Logger {
	logOutput := io.Writer(os.Stderr)
	var json bool
	var logLevel hclog.Level

	// if TF_LOG_PATH is set, output logs there
	if logPath := os.Getenv(envLogFile); logPath != "" {
		f, err := os.OpenFile(logPath, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		} else {
			logOutput = f
		}
	}

	// if TF_LOG is set, set the level
	envLevel := strings.ToUpper(os.Getenv(envLog))
	if envLevel == "" {
		logLevel = hclog.Off
	}
	if envLevel == "JSON" {
		logLevel = hclog.Trace
		json = true
	} else if isValidLogLevel(envLevel) {
		logLevel = hclog.LevelFromString(envLevel)
	} else {
		fmt.Fprintf(os.Stderr, "[WARN] Invalid log level: %q. Defaulting to level: OFF. Valid levels are: %+v",
			envLevel, ValidLevels)
	}

	return hclog.New(&hclog.LoggerOptions{
		Level:             logLevel,
		Output:            logOutput,
		IndependentLevels: true,
		JSONFormat:        json,
	})
}

func isValidLogLevel(level string) bool {
	for _, l := range ValidLevels {
		if level == string(l) {
			return true
		}
	}

	return false
}
