package scripts

func newConsoleLog(logger Logger) map[string]any {
	if logger == nil {
		return map[string]any{
			"log": func(args ...any) {},
		}
	}
	return map[string]any{
		"log": logger.Log,
	}
}
