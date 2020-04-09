package requestenvelopego

import (
	"context"

	tracergo "github.com/AccelByte/tracer-go"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func ChildScopeFromRemoteScope(
	rootCtx context.Context,
	rootLogger *logrus.Entry,
	name string,
	spanContextStr string,
	traceID string,
) *Scope {
	span, ctx := tracergo.ChildSpanFromRemoteSpan(rootCtx, name, spanContextStr)

	logger := rootLogger.WithFields(logrus.Fields{
		"caller": name,
		"trace":  traceID,
		"span":   tracergo.GetSpanContextString(span),
	})

	return &Scope{
		Ctx:     ctx,
		TraceID: traceID,
		Span:    span,
		Logger:  logger,
	}
}

func NewRootScope(rootCtx context.Context, rootLogger *logrus.Entry, name string, abTraceID string) *Scope {
	span, ctx := tracergo.StartSpanFromContext(rootCtx, name)

	logger := rootLogger.WithFields(logrus.Fields{
		"caller": name,
		"trace":  abTraceID,
		"span":   tracergo.GetSpanContextString(span),
	})

	scope := &Scope{
		Ctx:     ctx,
		TraceID: abTraceID,
		Span:    span,
		Logger:  logger,
	}

	if abTraceID != "" {
		scope.TraceTag(tracergo.TraceIDKey, abTraceID)
	}

	return scope
}

// Scope used as the envelope to combine and transport request-related information by the chain of function calls
type Scope struct {
	Ctx     context.Context
	TraceID string
	Span    opentracing.Span
	Logger  *logrus.Entry
}

// Finish finishes current scope
func (s *Scope) Finish() {
	tracergo.Finish(s.Span)
}

// TraceLog sends a log into tracer
func (s *Scope) TraceLog(key, value string) {
	tracergo.AddLog(s.Span, key, value)
}

// TraceLog sends a log into tracer
func (s *Scope) TraceError(err error) {
	tracergo.TraceError(s.Span, err)
}

// TraceTag sends a tag into tracer
func (s *Scope) TraceTag(key, value string) {
	tracergo.AddTag(s.Span, key, value)
}

// AddBaggage sends a baggage item into tracer
func (s *Scope) AddBaggage(key string, value string) {
	tracergo.AddBaggage(s.Span, key, value)
}

// GetSpanContextString gets scope span context string
func (s *Scope) GetSpanContextString() string {
	return tracergo.GetSpanContextString(s.Span)
}

// NewChildScope creates new child Scope
func (s *Scope) NewChildScope(name string) *Scope {
	span := opentracing.StartSpan(
		name,
		opentracing.ChildOf(s.Span.Context()),
	)

	logger := s.Logger.WithFields(logrus.Fields{
		"caller": name,
		"trace":  s.TraceID,
		"span":   tracergo.GetSpanContextString(span),
	})

	return &Scope{
		Ctx:     s.Ctx,
		TraceID: s.TraceID,
		Span:    span,
		Logger:  logger,
	}
}
