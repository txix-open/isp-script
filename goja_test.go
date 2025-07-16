package scripts_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	scripts "github.com/txix-open/isp-script"
)

func TestGoja_AddFunction(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	var response = {};
	try {
		response["test"] = f("string", "unknown", "test");
	} catch (e) {
		if (!(e instanceof GoError)) {
			throw(e);
		}
		if (e.value.Error() !== "TEST") {
			throw("Unexpected value: " + e.value.Error());
		}
	}
	return response;
	`

	f := func(varchar string, integer int, object string) (interface{}, error) {
		a.Equal("string", varchar)
		a.Equal(0, integer)
		a.Equal("test", object)
		return "test", nil
	}
	vm := goja.New()
	vm.Set("f", f)
	resp, err := vm.RunString(fmt.Sprintf("(function() { %s })();", SCRIPT))
	a.NoError(err)
	a.Equal(resp.Export(), map[string]interface{}{"test": "test"})

	f2 := func(varchar string, integer int, object string) (interface{}, error) {
		a.Equal("string", varchar)
		a.Equal(0, integer)
		a.Equal("test", object)
		return "test", errors.New("TEST")
	}
	vm = goja.New()
	vm.Set("f", f2)
	resp, err = vm.RunString(fmt.Sprintf("(function() { %s })();", SCRIPT))
	a.NoError(err)
	a.Equal(resp.Export(), map[string]interface{}{})
}

func TestScript_Default(t *testing.T) {
	a := assert.New(t)

	const SHARED = `
	function shared() {
		return arg["key"] + 2
	}
`
	const SCRIPT = `
	return shared()
`
	script, err := scripts.NewScript([]byte(SHARED), []byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)

	arg := map[string]interface{}{"key": 3, "4": 7}

	result, err := scripts.NewEngine().Execute(script, arg)
	a.NoError(err)

	a.Equal(int64(5), result)
}

func TestScript_WithLogging(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	console.log(arg)
	console.log(1, 2, 3)
	console.log("test")
	return 5
`
	script, err := scripts.NewScript([]byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)

	arg := map[string]interface{}{"key": 3}

	result, err := scripts.NewEngine().Execute(script, arg, scripts.WithLogger(scripts.NewStdoutJsonLogger()))
	a.NoError(err)
	a.Equal(int64(5), result)
}

func TestScript_WithMainFuncTrace(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	const x = 1
	function test(a) {
		return a + x
	}

	function main(arg) {
		return test(arg)
	}
`

	script, err := scripts.NewScript([]byte(SCRIPT))
	a.NoError(err)

	arg := 3

	engine := scripts.NewEngine()

	result, err := engine.Execute(script, arg, scripts.WithTraceMain())
	a.NoError(err)
	a.Equal(int64(arg+1), result)

	arg = 5

	result, err = engine.Execute(script, arg, scripts.WithTraceMain())
	a.NoError(err)
	a.Equal(int64(arg+1), result)
}

func TestScript_WithData(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	return i + str + mp[3]+arg["key"]+arr
`
	script, err := scripts.NewScript([]byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)

	arg := map[string]interface{}{"key": 3}

	result, err := scripts.NewEngine().Execute(script, arg,
		scripts.WithSet("i", 1),
		scripts.WithSet("str", "two"),
		scripts.WithSet("mp", map[string]interface{}{"3": "four"}),
		scripts.WithSet("arr", []int{5, 6, 7}))
	a.NoError(err)

	a.Equal("1twofour35,6,7", result)
}

func TestScript_WithFunc(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	return sqrt(arg["key"])
`
	script, err := scripts.NewScript([]byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)

	arg := map[string]interface{}{"key": 3}

	sqrt := func(x int) int {
		return x * x
	}
	result, err := scripts.NewEngine().Execute(script, arg, scripts.WithSet("sqrt", sqrt))
	a.NoError(err)

	a.Equal(int64(9), result)
}

func TestScript_WithDataWithFunc(t *testing.T) {
	a := assert.New(t)

	const SCRIPT = `
	return sqrt(arg["key"]) + sqrt(i)
`
	script, err := scripts.NewScript([]byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)

	arg := map[string]interface{}{"key": 3}

	sqrt := func(x int) int {
		return x * x
	}
	result, err := scripts.NewEngine().Execute(script, arg, scripts.WithSet("sqrt", sqrt), scripts.WithSet("i", 1))
	a.NoError(err)

	a.Equal(int64(10), result)
}
