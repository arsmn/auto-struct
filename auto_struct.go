package autostruct

import "reflect"

const defaultTag = "auto"

type option func(*config)

type config struct {
	tag string
}

func newConfig(opts ...option) *config {
	cfg := &config{
		tag: defaultTag,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func WithTag(tag string) option {
	return func(c *config) {
		c.tag = tag
	}
}

func Set(v any, opts ...option) error {
	return structFieldsSetter(newConfig(opts...), reflect.ValueOf(v))
}

func MustSet(v any, opts ...option) {
	if err := Set(v, opts...); err != nil {
		panic(err)
	}
}

func New[T any](opts ...option) T {
	var v T
	MustSet(&v, opts...)
	return v
}
