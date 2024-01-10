package scripts

import (
	"bytes"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

type Script struct {
	prog *goja.Program
}

func NewScript(source ...[]byte) (Script, error) {
	prog, err := goja.Compile("script", string(bytes.Join(source, []byte("\n"))), false)
	return Script{prog: prog}, err
}

type Engine struct {
	pool *sync.Pool
}

func NewEngine(opts ...EngineOption) *Engine {
	options := &engineOptions{
		moduleLoader: NewStaticModuleLoader(),
	}
	for _, opt := range opts {
		opt(options)
	}
	registry := require.NewRegistryWithLoader(options.moduleLoader.SourceLoader())
	return &Engine{
		pool: &sync.Pool{
			New: func() any {
				vm := goja.New()
				registry.Enable(vm)
				return vm
			},
		}}
}

func (m *Engine) Execute(s Script, arg any, opts ...ExecOption) (any, error) {
	options := &execOptions{
		arg:           arg,
		scriptTimeout: 1 * time.Second,
	}
	for _, o := range opts {
		o(options)
	}

	vm := m.pool.Get().(*goja.Runtime)
	vm.ClearInterrupt()
	options.set(vm)
	timer := time.AfterFunc(options.scriptTimeout, func() {
		vm.Interrupt("execution timeout")
	})
	defer func() {
		timer.Stop()
		vm.ClearInterrupt()
		options.reset(vm)
		m.pool.Put(vm)
	}()

	res, err := vm.RunProgram(s.prog)
	if err != nil {
		return nil, castErr(err)
	}

	return res.Export(), nil
}

func castErr(err error) error {
	if exception, ok := err.(*goja.Exception); ok {
		val := exception.Value().Export()
		if castedErr, ok := val.(error); ok {
			return castedErr
		}
	}
	return err
}
