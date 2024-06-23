package supervisor

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wastedcode/supervisor/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Child struct {
	logger *slog.Logger
	span   trace.Span
}

func (c Child) Log(ctx context.Context, msg string, attrs ...slog.Attr) {
	c.logger.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func (c Child) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	c.logger.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

func (c Child) End() {
	c.span.End()
}

func (c Child) Error(err error) error {
	if err == nil {
		return nil
	}

	c.span.RecordError(err)
	c.span.SetStatus(codes.Error, fmt.Sprintf("%s", err))
	c.logger.Error(errors.FlattenDetails(err), slog.Any("err", newErrLog(err)))
	return err
}
