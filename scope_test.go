package requestenvelopego

import (
	"context"
	"testing"

	tracergo "github.com/AccelByte/tracer-go"
	"github.com/stretchr/testify/require"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	jaegerclientgo "github.com/uber/jaeger-client-go"
)

func TestGetSpanContextString(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	closer := tracergo.InitGlobalTracer("", "", "test", "")
	defer closer.Close()
}

func TestChildScopeFromRemoteScope(t *testing.T) {
	logger := logrus.WithField("test", "test")

	closer := tracergo.InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	expectedSpan, _ := opentracing.StartSpanFromContext(context.Background(), "test")

	spanContextStr := expectedSpan.Context().(jaegerclientgo.SpanContext).String()

	scope := ChildScopeFromRemoteScope(context.Background(), logger, "test", spanContextStr, "test-trace-id")

	scope.Logger.Errorf("")

	assert.Equal(t,
		expectedSpan.Context().(jaegerclientgo.SpanContext).TraceID().String(),
		scope.Span.Context().(jaegerclientgo.SpanContext).TraceID().String(),
	)

	assert.Equal(t,
		expectedSpan.Context().(jaegerclientgo.SpanContext).SpanID().String(),
		scope.Span.Context().(jaegerclientgo.SpanContext).ParentID().String(),
	)
}

func TestChildScopeFromRemoteScope_EmptySpanContextString(t *testing.T) {
	logger := logrus.WithField("test", "test")

	closer := tracergo.InitGlobalTracer("", "", "test", "")
	defer closer.Close()

	scope := ChildScopeFromRemoteScope(context.Background(), logger, "test", "", "test-trace-id")

	scope.Logger.Println("")

	assert.NotEmpty(t,
		scope.Span.Context().(jaegerclientgo.SpanContext).TraceID().String(),
	)

	assert.NotEmpty(t,
		scope.Span.Context().(jaegerclientgo.SpanContext).ParentID().String(),
	)
}

func TestNewRootScope(t *testing.T) {
	logger := logrus.WithField("test", "test")

	scope := NewRootScope(context.Background(), logger, "name", "trace-id")
	defer scope.Finish()

	require.NotNil(t, scope)
}

func TestScope_NewChildScope(t *testing.T) {
	logger := logrus.WithField("test", "test")

	scope := NewRootScope(context.Background(), logger, "name", "trace-id")
	defer scope.Finish()

	require.NotNil(t, scope)

	childScope := scope.NewChildScope("name2")
	defer childScope.Finish()

	require.NotNil(t, childScope)
}
