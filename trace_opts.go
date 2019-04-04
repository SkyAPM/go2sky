package go2sky

// WithReporter setup report pipeline for tracer
func WithReporter(reporter Reporter) TracerOption {
	return func(t *Tracer) {
		t.reporter = reporter
	}
}
