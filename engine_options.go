package scripts

import (
	"github.com/dop251/goja_nodejs/require"
)

type engineOptions struct {
	moduleLoader ModuleLoader
	pathResolver require.PathResolver
}

type EngineOption func(opts *engineOptions)

func WithModuleLoader(moduleLoader ModuleLoader) EngineOption {
	return func(opts *engineOptions) {
		opts.moduleLoader = moduleLoader
	}
}

func WithPathResolver(pathResolver require.PathResolver) EngineOption {
	return func(opts *engineOptions) {
		opts.pathResolver = pathResolver
	}
}
