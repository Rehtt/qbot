package cqhttp_bot

type Options struct {
	handleThreadNum int
	debug           bool
}
type Option func(options *Options)

// WithHandleThreadNum 默认为200，如果数值太小会导致处理出现延时
func WithHandleThreadNum(n int) Option {
	return func(options *Options) {
		options.handleThreadNum = n
	}
}
func WithDebug(n bool) Option {
	return func(options *Options) {
		options.debug = n
	}
}

func loadOptions(options ...Option) *Options {
	o := new(Options)
	for _, opt := range options {
		opt(o)
	}
	return o
}
