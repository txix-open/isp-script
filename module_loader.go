package scripts

import (
	"fmt"
	"io"
	"io/fs"
)

type ModuleLoader interface {
	SourceLoader() func(path string) ([]byte, error)
}

type StaticModule struct {
	name string
	data []byte
}

func NewStaticModule(name string, data []byte) StaticModule {
	return StaticModule{
		name: name,
		data: data,
	}
}

type StaticModuleLoader struct {
	modules map[string][]byte
}

func NewStaticModuleLoader(staticModules ...StaticModule) StaticModuleLoader {
	modules := map[string][]byte{}
	for _, module := range staticModules {
		modules[module.name] = module.data
	}
	return StaticModuleLoader{
		modules: modules,
	}
}

func (m StaticModuleLoader) SourceLoader() func(path string) ([]byte, error) {
	return func(path string) ([]byte, error) {
		module, ok := m.modules[path]
		if !ok {
			return nil, fmt.Errorf("module %s does not exist", path)
		}
		return module, nil
	}
}

type FsModuleLoader struct {
	fs fs.FS
}

func NewFsModuleLoader(fs fs.FS) FsModuleLoader {
	return FsModuleLoader{
		fs: fs,
	}
}

func (m FsModuleLoader) SourceLoader() func(path string) ([]byte, error) {
	return func(path string) ([]byte, error) {
		file, err := m.fs.Open(path)
		if err != nil {
			return nil, fmt.Errorf("fs open: %w", err)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}

		return data, nil
	}
}
