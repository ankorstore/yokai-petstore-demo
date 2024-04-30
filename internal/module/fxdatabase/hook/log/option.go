package log

import (
	"github.com/ankorstore/yokai/log"
	"github.com/rs/zerolog"
)

type Options struct {
	Level              zerolog.Level
	Arguments          bool
	ExcludedOperations []string
}

func DefaultLogHookOptions() Options {
	return Options{
		Level:              zerolog.DebugLevel,
		Arguments:          false,
		ExcludedOperations: []string{},
	}
}

type LogHookOption func(o *Options)

func WithLevel(level string) LogHookOption {
	return func(o *Options) {
		o.Level = log.FetchLogLevel(level)
	}
}

func WithArguments(arguments bool) LogHookOption {
	return func(o *Options) {
		o.Arguments = arguments
	}
}

func WithExcludedOperations(excludedOperations ...string) LogHookOption {
	return func(o *Options) {
		o.ExcludedOperations = excludedOperations
	}
}
