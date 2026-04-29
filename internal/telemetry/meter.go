package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// initMeterProvider builds an OTel meter provider that exports via OTLP/HTTP
// on a periodic interval. Endpoint, headers, etc. are picked up from the
// standard OTEL_* env vars (same as the tracer and logger).
func initMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(provider)
	return provider, nil
}
