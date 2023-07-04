package scripts

import (
	"encoding/json"
	"os"
)

type Logger interface {
	Log(arg ...any)
}

type StdoutJsonLogger struct{}

func NewStdoutJsonLogger() StdoutJsonLogger {
	return StdoutJsonLogger{}
}

func (s StdoutJsonLogger) Log(args ...any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(args)
	_, _ = os.Stdout.WriteString("\n")
}
