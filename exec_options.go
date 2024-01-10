package scripts

import (
	"time"

	"github.com/dop251/goja"
)

type ExecOption func(opt *options)

type options struct {
	scriptTimeout   time.Duration
	arg             any
	logger          Logger
	data            map[string]any
	fieldNameMapper goja.FieldNameMapper
}

func WithTimeout(duration time.Duration) ExecOption {
	return func(opt *options) {
		opt.scriptTimeout = duration
	}
}

func WithLogger(logger Logger) ExecOption {
	return func(opt *options) {
		opt.logger = logger
	}
}

func WithSet(name string, f any) ExecOption {
	return func(opt *options) {
		if opt.data == nil {
			opt.data = make(map[string]any)
		}
		opt.data[name] = f
	}
}

func WithDefaultToolkit() ExecOption {
	return WithSet("toolkit", toolkit)
}

func WithFieldNameMapper(fieldNameMapper goja.FieldNameMapper) ExecOption {
	return func(opt *options) {
		opt.fieldNameMapper = fieldNameMapper
	}
}

func (c *options) set(vm *goja.Runtime) {
	if c.fieldNameMapper != nil {
		vm.SetFieldNameMapper(c.fieldNameMapper)
	}
	vm.Set("arg", c.arg)
	console := newConsoleLog(c.logger)
	vm.Set("console", console)
	for name, data := range c.data {
		vm.Set(name, data)
	}
}

func (c *options) reset(vm *goja.Runtime) {
	vm.SetFieldNameMapper(nil)
	vm.Set("arg", goja.Undefined())
	vm.Set("console", goja.Undefined())
	for name := range c.data {
		vm.Set(name, goja.Undefined())
	}
}
