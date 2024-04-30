package trace

type Options struct {
	Arguments          bool
	ExcludedOperations []string
}

func DefaultTraceHookOptions() Options {
	return Options{
		Arguments:          false,
		ExcludedOperations: []string{},
	}
}

type TraceHookOption func(o *Options)

func WithArguments(arguments bool) TraceHookOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}

func WithExcludedOperations(excludedOperations ...string) TraceHookOption {
	return func(o *Options) {
		o.ExcludedOperations = excludedOperations
	}
}
