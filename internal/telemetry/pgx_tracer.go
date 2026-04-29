package telemetry

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// pgxTracer implements pgx.QueryTracer to emit one span per query.
// pgx calls TraceQueryStart before sending the query and TraceQueryEnd
// after the result is read; the ctx returned from Start carries the span
// so End can pull it back out.
type pgxTracer struct {
	tracer trace.Tracer
}

// NewPgxTracer returns a pgx.QueryTracer that records each query as an OTel
// span. Wire it into pgxpool.Config.ConnConfig.Tracer.
func NewPgxTracer() pgx.QueryTracer {
	return &pgxTracer{tracer: otel.Tracer("github.com/urlspace/api/internal/telemetry")}
}

func (t *pgxTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx, _ = t.tracer.Start(ctx, queryName(data.SQL),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", data.SQL),
		),
	)
	return ctx
}

func (t *pgxTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if data.Err != nil {
		span.RecordError(data.Err)
		span.SetStatus(codes.Error, data.Err.Error())
		return
	}
	span.SetAttributes(attribute.String("db.command_tag", data.CommandTag.String()))
}

// queryName picks a human-readable span name. sqlc prepends each query with
// `-- name: <Name> :<kind>` and pgx forwards the whole string, so we prefer
// that name (e.g. GetUserByID) over the bare SQL keyword. Falls back to the
// first SQL keyword for raw queries (BEGIN, COMMIT, ad-hoc Exec/Query).
func queryName(sql string) string {
	sql = strings.TrimSpace(sql)
	if name := sqlcName(sql); name != "" {
		return name
	}
	if sql == "" {
		return "query"
	}
	if i := strings.IndexAny(sql, " \t\n"); i > 0 {
		return strings.ToUpper(sql[:i])
	}
	return strings.ToUpper(sql)
}

// sqlcName extracts "GetUserByID" from a leading "-- name: GetUserByID :one"
// sqlc comment. Returns "" if the SQL doesn't start with that marker.
func sqlcName(sql string) string {
	const prefix = "-- name:"
	if !strings.HasPrefix(sql, prefix) {
		return ""
	}
	rest := strings.TrimLeft(sql[len(prefix):], " \t")
	if i := strings.IndexAny(rest, " \t\n"); i > 0 {
		return rest[:i]
	}
	return rest
}
