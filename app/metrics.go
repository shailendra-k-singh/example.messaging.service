package app

import (
	"io"
	"net/http"
	"net/http/httptrace"

	"github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegerprom "github.com/uber/jaeger-lib/metrics/prometheus"

	log "github.com/sirupsen/logrus"
)

func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace {
	trace := &clientTrace{span: span}
	return &httptrace.ClientTrace{
		DNSStart: trace.dnsStart,
		DNSDone:  trace.dnsDone,
	}
}

// clientTrace holds a reference to the Span and
// provides methods used as ClientTrace callbacks
type clientTrace struct {
	span opentracing.Span
}

func (h *clientTrace) dnsStart(info httptrace.DNSStartInfo) {
	h.span.LogKV(
		otlog.String("event", "DNS start"),
		otlog.Object("host", info.Host),
	)
}

func (h *clientTrace) dnsDone(httptrace.DNSDoneInfo) {
	h.span.LogKV(otlog.String("event", "DNS done"))
}

type LogrusAdapter struct{}

func (l LogrusAdapter) Error(msg string) {
	log.Errorf(msg)
}

func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

func (l LogrusAdapter) Debugf(msg string, args ...interface{}) {
	log.Debugf(msg, args...)
}

type tracerObj struct {
	closer   io.Closer
	reporter jaeger.Reporter
	tracer   opentracing.Tracer
}

func (t *tracerObj) initTracing(lib string, host, address string) error {
	factory := jaegerprom.New()
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": lib})
	log.Info("setting jaeger agent as: ",host + address)
	transport, err := jaeger.NewUDPTransport(host + address, 0)
	if err != nil {
		return err
	}
	logAdapt := LogrusAdapter{}
	t.reporter = jaeger.NewCompositeReporter(
		jaeger.NewLoggingReporter(logAdapt),
		jaeger.NewRemoteReporter(transport,
			jaeger.ReporterOptions.Metrics(metrics),
			jaeger.ReporterOptions.Logger(logAdapt),
		),
	)
	sampler := jaeger.NewConstSampler(true)
	t.tracer, t.closer = jaeger.NewTracer("Messaging Service",
		sampler,
		t.reporter,
		jaeger.TracerOptions.Metrics(metrics),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)
	log.Debug("Tracing initialized")
	opentracing.SetGlobalTracer(t.tracer)
	return nil
}

func (t *tracerObj) startTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		span := t.tracer.StartSpan(req.URL.Path)
		span.SetTag(string(ext.Component), "client")
		defer span.Finish()

		ctx := opentracing.ContextWithSpan(req.Context(), span)
		// add http tracing
		trace := NewClientTrace(span)
		ctx = httptrace.WithClientTrace(ctx, trace)
		req = req.WithContext(ctx)

		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			log.Errorf("error while injecting context is: %v", err)
		}
		next.ServeHTTP(w, req)
	})
}

func (t *tracerObj) Close() {
	t.reporter.Close()
	err := t.closer.Close()
	if err != nil {
		log.Error("Error while closing closer object: ", err)
	}
}
