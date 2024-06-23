package supervisor

import "log/slog"

type Config struct {
	LogLevel        slog.Level
	ApplicationName string
	OtelEnvironment string
}
