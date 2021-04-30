package tflog

type loggerKey string

var (
	providerSpaceRootLoggerKey loggerKey = "provider"
	sdkRootLoggerKey           loggerKey = "sdk"
)
