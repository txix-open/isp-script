package scripts_test

import (
	"embed"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	scripts "github.com/txix-open/isp-script"
)

//go:embed test_modules
var fs embed.FS

func TestEmbedFsModules(t *testing.T) {
	require := require.New(t)

	engine := scripts.NewEngine(
		scripts.WithModuleLoader(
			scripts.NewFsModuleLoader(fs),
		),
	)
	src, err := scripts.NewScript([]byte(`
	const module = require("./test_modules/a.js");
	module.a();
`))
	require.NoError(err)
	result, err := engine.Execute(src, nil)
	require.NoError(err)
	require.EqualValues("hello world", result)
}

func TestStaticModules(t *testing.T) {
	require := require.New(t)

	moduleData, err := os.ReadFile("test_modules/b.js")
	require.NoError(err)

	engine := scripts.NewEngine(
		scripts.WithModuleLoader(
			scripts.NewStaticModuleLoader(
				scripts.NewStaticModule("b.js", moduleData),
			),
		),
	)
	src, err := scripts.NewScript([]byte(`(function() {
	const module = require("./b.js");
	return module.b();
})();
`))
	require.NoError(err)
	result, err := engine.Execute(src, nil, scripts.WithTimeout(100*time.Second))
	require.NoError(err)
	require.EqualValues("hello world", result)

	result, err = engine.Execute(src, nil, scripts.WithTimeout(100*time.Second))
	require.NoError(err)
	require.EqualValues("hello world", result)
}
