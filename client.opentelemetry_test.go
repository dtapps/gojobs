package gojobs

import (
	"context"
	"testing"
)

func TestTrace(t *testing.T) {
	ctx, span := TraceStartSpan(context.TODO(), "TestTrace")
	t.Log(TraceGetTraceID(ctx))
	t.Log(TraceGetSpanID(ctx))
	defer TraceEndSpan(span)
}

func TestTraceStartSpan(t *testing.T) {
	TraceStartSpan(context.TODO(), "TestTraceStartSpan")
	TraceStartSpan(context.Background(), "TestTraceStartSpan")
}

func TestTraceGetTraceID(t *testing.T) {
	t.Log(TraceGetTraceID(context.TODO()))
	t.Log(TraceGetTraceID(context.Background()))
}

func TestTraceGetSpanID(t *testing.T) {
	t.Log(TraceGetSpanID(context.TODO()))
	t.Log(TraceGetSpanID(context.Background()))
}
