package go2sky

func WithReporter(reporter Reporter) TracerOption{
	return func(t *Tracer) {
		t.reporter = reporter
	}
}
