package tracing

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func getQueryName(query string) (name string) {
	lines := strings.SplitN(query, "\n", 1)
	name = "pgx"
	if len(lines) > 0 {
		first := lines[0]
		params := strings.SplitN(first, ":", 3)
		if len(params) > 1 {
			name += "." + strings.TrimSpace(params[1])
		}
	}

	return name
}

func SpanFromContext(ctx context.Context, name string, opts ...tracer.StartSpanOption) (ddtrace.Span, context.Context) {
	parent, ok := tracer.SpanFromContext(ctx)
	if !ok {
		return tracer.StartSpanFromContext(ctx, name, opts...)
	}

	optsWithChild := []tracer.StartSpanOption{tracer.ChildOf(parent.Context())}
	optsWithChild = append(optsWithChild, opts...)
	span := tracer.StartSpan(
		name,
		optsWithChild...,
	)

	return span, tracer.ContextWithSpan(ctx, span)
}

type DBTracer struct {
	serviceName string
}

func NewDBTracer(serviceName string) *DBTracer {
	return &DBTracer{
		serviceName: serviceName,
	}
}

func (l *DBTracer) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]any) {
	t, timeOk := data["time"].(time.Duration)
	q, queryOk := data["sql"].(string)
	if timeOk && queryOk {
		end := time.Now()
		start := end.Add(-t)
		startTime := tracer.StartTime(start)

		opts := []ddtrace.StartSpanOption{
			tracer.SpanType(ext.SpanTypeSQL),
			tracer.ResourceName(string(q)),
			tracer.ServiceName(l.serviceName),
			startTime,
		}

		span, _ := SpanFromContext(ctx, getQueryName(q), opts...)
		span.Finish(tracer.FinishTime(end))
	}
}
