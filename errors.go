package supervisor

import "github.com/wastedcode/supervisor/errors"

var (
	ErrUnsupportedOtelEnv = errors.New("unsupported otel environment")
)
