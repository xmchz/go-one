package core

type Option func(*logger)

func WithWriters(writers ...Writer) Option {
	return func(l *logger) {
		l.writers = writers
	}
}

func WithLevel(level Level) Option {
	return func(l *logger) {
		l.level = level
	}
}

func WithCallerSkip(skip int) Option {
	return func(l *logger) {
		l.callerSkip = skip
	}
}

