package go_fsm

// Logger is a generic logging interface
type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

// Nil logger adapter (logging nothing)
type nilLoggerAdapter struct {
}

func (n *nilLoggerAdapter) Log(v ...interface{}) {
}

func (n *nilLoggerAdapter) Logf(format string, v ...interface{}) {
}
