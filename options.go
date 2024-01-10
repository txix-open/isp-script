package scripts

type engineOptions struct {
	moduleLoader ModuleLoader
}

type EngineOption func(opts *engineOptions)

func WithModuleLoader(moduleLoader ModuleLoader) EngineOption {
	return func(opts *engineOptions) {
		opts.moduleLoader = moduleLoader
	}
}
