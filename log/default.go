package log

import "os"

const (
	defaultLoggerCh = 20000
)

var (
	defaultLogger *logger
)

func Init(opts ...Option) {
	defaultLogger = &logger{
		dataCh:     make(chan *Data, defaultLoggerCh),
		callerSkip: 2,
		level:      InfoLevel,
	}
	for _, opt := range opts {
		opt(defaultLogger)
	}
	defaultLogger.wg.Add(1)
	go defaultLogger.run()
}

func Stop() {
	close(defaultLogger.dataCh)
	defaultLogger.wg.Wait()
	for _, writer := range defaultLogger.writers {
		writer.Close()
	}
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Log(DebugLevel, format, args...)
}

func Trace(format string, args ...interface{}) {
	defaultLogger.Log(TraceLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Log(InfoLevel, format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Log(WarnLevel, format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Log(ErrorLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Log(ErrorLevel, format, args...)
	os.Exit(1)
}
