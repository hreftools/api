package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Setup installs the slog default and wires up OTel tracing + logging.
// First it sets a stdout-only slog handler so early init errors are visible,
// then it brings up the tracer and logger providers and upgrades slog to fan
// out to both stdout and the OTel bridge.
//
// The returned shutdown function flushes both providers; callers should defer
// it with a fresh context (the parent ctx may already be cancelled at exit).
func Setup(ctx context.Context) (shutdown func(context.Context) error, err error) {
	slog.SetDefault(slog.New(newStdoutHandler()))

	tp, err := initTracerProvider(ctx)
	if err != nil {
		return nil, fmt.Errorf("tracer provider: %w", err)
	}

	lp, err := initLoggerProvider(ctx)
	if err != nil {
		_ = tp.Shutdown(ctx)
		return nil, fmt.Errorf("logger provider: %w", err)
	}

	attachOtelLogger(lp)

	return func(ctx context.Context) error {
		return errors.Join(tp.Shutdown(ctx), lp.Shutdown(ctx))
	}, nil
}
