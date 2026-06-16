package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	Ctx           context.Context
	CollectorHost string
	CollectorPort string
	ServiceName   string
}

func New(
	ctx context.Context,
	collectorHost string,
	collectorPort string,
	serviceName string,
) *Telemetry {
	return &Telemetry{
		Ctx:           ctx,
		CollectorHost: collectorHost,
		CollectorPort: collectorPort,
		ServiceName:   serviceName,
	}
}

func (Telemetry) propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (t *Telemetry) traceProvider() (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(t.Ctx,
		otlptracegrpc.WithEndpoint(
			fmt.Sprintf("%s:%s", t.CollectorHost, t.CollectorPort),
		),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := t.resource()
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(
			traceExporter,
			sdktrace.WithBatchTimeout(time.Second),
		),
	)

	return tracerProvider, nil
}

func (t *Telemetry) logProvider() (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(t.Ctx,
		otlploggrpc.WithEndpoint(
			fmt.Sprintf("%s:%s", t.CollectorHost, t.CollectorPort),
		),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := t.resource()
	if err != nil {
		return nil, err
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(
				logExporter,
				sdklog.WithExportInterval(time.Second),
			),
		),
	)

	return loggerProvider, nil
}

func (t *Telemetry) resource() (*resource.Resource, error) {
	res, err := resource.New(t.Ctx,
		resource.WithAttributes(
			semconv.ServiceName(t.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *Telemetry) NewTracer() (trace.Tracer, error) {
	otel.SetTextMapPropagator(t.propagator())

	tp, err := t.traceProvider()
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)

	lp, err := t.logProvider()
	if err != nil {
		return nil, err
	}
	global.SetLoggerProvider(lp)

	return tp.Tracer(t.ServiceName), nil
}
