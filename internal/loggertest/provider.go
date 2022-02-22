package loggertest

import (
	"context"
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

func ProviderRoot(ctx context.Context, output io.Writer) context.Context {
	loggerOptions := &hclog.LoggerOptions{
		DisableTime: true,
		JSONFormat:  true,
		Level:       hclog.Trace,
		Output:      output,
	}

	ctx = logging.SetProviderRootLogger(ctx, hclog.New(loggerOptions))

	return ctx
}
