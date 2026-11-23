package scripts

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
)

var ErrScriptNotCompiled = errors.New("script is not compiled")

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
		pathResolver: require.DefaultPathResolver,
	}
	for _, opt := range opts {
		opt(options)
	}
	registry := require.NewRegistry(
		require.WithLoader(options.moduleLoader.SourceLoader()),
		require.WithPathResolver(options.pathResolver),
	)
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
	return m.withVM(arg, opts, func(vm *goja.Runtime) (any, error) {
		if s.prog == nil {
			return nil, ErrScriptNotCompiled
		}
		res, err := vm.RunProgram(s.prog)
		if err != nil {
			return nil, castErr(err)
		}

		return res.Export(), nil
	})
}

// ExecutePipeline последовательно исполняет скрипты в одной VM,
// результат каждого скрипта становится arg следующего
func (m *Engine) ExecutePipeline(scripts []Script, arg any, opts ...ExecOption) (any, error) {
	return m.withVM(arg, opts, func(vm *goja.Runtime) (any, error) {
		if len(scripts) == 0 {
			return nil, nil
		}

		var lastResult any
		for i, script := range scripts {
			if script.prog == nil {
				return nil, ErrScriptNotCompiled
			}
			res, err := vm.RunProgram(script.prog)
			if err != nil {
				return nil, castErr(err)
			}

			lastResult = res.Export()
			if lastResult == nil {
				return nil, nil
			}

			if i < len(scripts)-1 {
				err = vm.Set("arg", res)
				if err != nil {
					return nil, castErr(err)
				}
			}
		}

		return lastResult, nil
	})
}

// ExecuteList исполняет все скрипты в одной VM,
// каждый скрипт работает с исходным arg
func (m *Engine) ExecuteList(scripts []Script, arg any, opts ...ExecOption) ([]any, error) {
	res, err := m.withVM(arg, opts, func(vm *goja.Runtime) (any, error) {
		if len(scripts) == 0 {
			return nil, nil
		}

		results := make([]any, 0, len(scripts))
		for _, script := range scripts {
			if script.prog == nil {
				return nil, ErrScriptNotCompiled
			}
			res, err := vm.RunProgram(script.prog)
			if err != nil {
				return nil, castErr(err)
			}

			results = append(results, res.Export())
		}

		return results, nil
	})
	if err != nil {
		return nil, err
	}

	arrRes, ok := res.([]any)
	if !ok {
		return nil, nil
	}
	return arrRes, nil
}

// ExecuteListParallel исполняет все скрипты параллельно,
// каждый скрипт в отдельной VM, каждый скрипт работает с исходным arg.
// arg должен быть read-only: один и тот же Go-объект используется всеми скриптами.
func (m *Engine) ExecuteListParallel(scripts []Script, arg any, opts ...ExecOption) ([]any, error) {
	if len(scripts) == 0 {
		return nil, nil
	}

	type scriptResponse struct {
		index int
		data  any
		err   error
	}
	ch := make(chan scriptResponse, len(scripts))
	var wg sync.WaitGroup

	for i, script := range scripts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := m.Execute(script, arg, opts...)
			ch <- scriptResponse{
				index: i,
				data:  res,
				err:   err,
			}
		}()
	}

	wg.Wait()
	close(ch)

	results := make([]any, len(scripts))

	for resp := range ch {
		if resp.err != nil {
			return nil, resp.err
		}

		results[resp.index] = resp.data
	}

	return results, nil
}

func (m *Engine) withVM(arg any, opts []ExecOption, fn func(vm *goja.Runtime) (any, error)) (any, error) {
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

	return fn(vm)
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
