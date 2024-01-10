package scripts_test

import (
	"fmt"
	"testing"

	scripts "github.com/integration-system/isp-script"
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
    const jsDate = toolkit.goTimeToDate(future);
    console.log(jsDate.getTime());
	return 5
`
	script, err := scripts.NewScript([]byte(fmt.Sprintf("(function() { %s })();", SCRIPT)))
	a.NoError(err)
	result, err := scripts.NewEngine().Execute(
		script,
		"2006-01-02 15:04:05",
		scripts.WithDefaultToolkit(),
		scripts.WithLogger(scripts.NewStdoutJsonLogger()),
	)
	a.NoError(err)
	a.Equal(int64(5), result)
}
