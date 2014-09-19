package level

type Level int

const (
	Off Level = iota
	Trace
	Debug
	Info
	Warn
	Error
	Fatal
	All
)

func (l Level) String() string {
	switch l {
	case Off:
		return "OFF"
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	case All:
		return "ALL"
	}

	return "Level-not-found"
}
