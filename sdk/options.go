package tfsdklog

import (
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type Option func(loggerOpts) loggerOpts

type loggerOpts struct {
	name            string
	level           hclog.Level
	includeLocation bool

	// some private options to be used only by tflog for testing purposes
	// we should never export an Option that sets these
	output io.Writer
}

func applyLoggerOpts(opts ...Option) loggerOpts {
	// set some defaults
	l := loggerOpts{
		level:           hclog.Trace,
		includeLocation: true,
	}
	for _, opt := range opts {
		l = opt(l)
	}
	return l
}

func withOutput(output io.Writer) Option {
	return func(l loggerOpts) loggerOpts {
		l.output = output
		return l
	}
}

func WithLogName(name string) Option {
	return func(l loggerOpts) loggerOpts {
		l.name = name
		return l
	}
}

func WithLevelFromEnv(name string, subsystems ...string) Option {
	return func(l loggerOpts) loggerOpts {
		envVar := strings.Join(subsystems, "_")
		if envVar != "" {
			envVar = "_" + envVar
		}
		envVar = strings.ToUpper(name + envVar)
		l.level = hclog.LevelFromString(os.Getenv(envVar))
		return l
	}
}

func WithoutLocation() Option {
	return func(l loggerOpts) loggerOpts {
		l.includeLocation = false
		return l
	}
}
