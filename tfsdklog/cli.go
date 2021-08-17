package tfsdklog

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

const envLogCLI = "TF_LOG_CLI"

func newCLILogger(ctx context.Context, commandID string) hclog.Logger {
	sink := getSink(ctx)
	if sink == nil {
		return nil
	}
	l := sink.Named("terraform")
	envLevel := strings.ToUpper(os.Getenv(envLog))
	if envLevel != "" {
		if isValidLogLevel(envLevel) {
			l.SetLevel(hclog.LevelFromString(envLevel))
		} else {
			fmt.Fprintf(os.Stderr, "[WARN] Invalid log level for %s: %q. Defaulting to %s level. Valid levels are: %+v",
				envLogCLI, envLevel, envLog, ValidLevels)
		}
	}
	l = l.With("command_id", commandID)
	return l
}
