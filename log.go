package supervisor

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/wastedcode/supervisor/errors"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

const (
	KeyRequestID              = "request.id"
	KeyHttpRoute              = string(semconv.HTTPRouteKey)
	KeyHttpRequest            = "http.request"
	KeyHttpMethod             = string(semconv.HTTPRequestMethodKey)
	KeyHttpPath               = "http.request.path"
	KeyHttpRequestClientIP    = string(semconv.ClientAddressKey)
	KeyHttpResponseSize       = string(semconv.HTTPResponseSizeKey)
	KeyHttpRequestDuration    = string(semconv.HTTPClientRequestDurationName)
	KeyHttpResponseStatusCode = string(semconv.HTTPStatusCodeKey)
)

func getSimpleJsonLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

type HttpRequest struct {
	r *http.Request
}

func NewHttpRequest(r *http.Request) *HttpRequest {
	return &HttpRequest{r: r}
}

func (r HttpRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String(KeyHttpMethod, r.r.Method),
		slog.String(KeyHttpPath, r.r.URL.Path),
		slog.String(KeyHttpRequestClientIP, GetIPFromHeaders(r.r).String()),
	)
}

type errLog struct {
	err error
}

func newErrLog(err error) errLog {
	return errLog{err: err}
}

func (err errLog) LogValue() slog.Value {
	return errors.LogValue(err.err)
}
