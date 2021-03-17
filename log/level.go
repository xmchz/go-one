package log

type Level int8

const (
	_ Level = iota
	DebugLevel
	TraceLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

//func (l Level) MarshalJSON() ([]byte, error) {
//	return []byte(l.String()), nil
//}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "TRACE"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	}
	return "UNKNOWN"
}
