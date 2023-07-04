package scripts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScript_Toolkit(t *testing.T) {
	a := require.New(t)

	const SCRIPT = `
	console.log(arg);
	console.log(toolkit.sha256(arg));
	console.log(toolkit.sha512(arg));
    console.log(toolkit.uuid());
    const date = toolkit.parseTime(arg, "15:04:05");
    console.log(date);
    const future = toolkit.now().Add(toolkit.durationFromMillis(1000));
    console.log(date.Before(future));
    console.log(date.Format("2006-01-02"));
	return 5
`
	script, err := NewScript([]byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)
	result, err := NewEngine().Execute(
		script,
		"2006-01-02 15:04:05",
		WithDefaultToolkit(),
		WithLogger(NewStdoutJsonLogger()),
	)
	a.NoError(err)
	a.Equal(int64(5), result)
}
