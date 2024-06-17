package errors

import (
	"fmt"
	"log/slog"

	pkgerr "github.com/cockroachdb/errors"
)

func New(msg string, args ...any) error {
	return pkgerr.Newf(msg)
}

func WithStack(err error) error {
	return pkgerr.WithStack(err)
}

func WithExternalDetails(err error, msg string) error {
	return pkgerr.WithHint(err, msg)
}

func WithInternalDetails(err error, msg string, args ...any) error {
	return pkgerr.WithDetailf(err, msg, args...)
}

func WithInternalDetailsAndStack(err error, msg string, args ...any) error {
	return WithInternalDetails(WithStack(err), msg, args...)
}

func Join(errs ...error) error {
	return pkgerr.Join(errs...)
}

func FlattenDetails(err error) string {
	return pkgerr.FlattenDetails(err)
}

func LogValue(err error) slog.Value {
	return slog.GroupValue(
		slog.String("error.message", pkgerr.FlattenDetails(err)),
		slog.Any("error.hints", pkgerr.GetAllDetails(err)),
		slog.Any("error.stack", fmt.Sprintf("%+v\n", err)),
	)
}
