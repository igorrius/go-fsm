package go_fsm

type Option func(*Options)

type Options struct {
	Logger Logger
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Logger: &nilLoggerAdapter{},
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func LoggerOption(l Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
