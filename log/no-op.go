package log

type NoopLogger struct{}

func (n NoopLogger) Debug(_ ...any)    {}
func (n NoopLogger) Info(_ ...any)     {}
func (n NoopLogger) Warning(_ ...any)  {}
func (n NoopLogger) Error(_ ...any)    {}
func (n NoopLogger) Critical(_ ...any) {}
func (n NoopLogger) Fatal(_ ...any)    {}
