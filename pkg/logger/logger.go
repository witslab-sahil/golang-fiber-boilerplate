package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func New(level string) Logger {
	log := logrus.New()
	
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	
	log.SetOutput(os.Stdout)
	
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
	
	return &logrusLogger{
		logger: log,
		entry:  nil,
	}
}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.entry != nil {
		l.entry.Debug(args...)
	} else {
		l.logger.Debug(args...)
	}
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Debugf(format, args...)
	} else {
		l.logger.Debugf(format, args...)
	}
}

func (l *logrusLogger) Info(args ...interface{}) {
	if l.entry != nil {
		l.entry.Info(args...)
	} else {
		l.logger.Info(args...)
	}
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Infof(format, args...)
	} else {
		l.logger.Infof(format, args...)
	}
}

func (l *logrusLogger) Warn(args ...interface{}) {
	if l.entry != nil {
		l.entry.Warn(args...)
	} else {
		l.logger.Warn(args...)
	}
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Warnf(format, args...)
	} else {
		l.logger.Warnf(format, args...)
	}
}

func (l *logrusLogger) Error(args ...interface{}) {
	if l.entry != nil {
		l.entry.Error(args...)
	} else {
		l.logger.Error(args...)
	}
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Errorf(format, args...)
	} else {
		l.logger.Errorf(format, args...)
	}
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatal(args...)
	} else {
		l.logger.Fatal(args...)
	}
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	if l.entry != nil {
		l.entry.Fatalf(format, args...)
	} else {
		l.logger.Fatalf(format, args...)
	}
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.logger.WithFields(fields),
	}
}

// GetLogrusLogger returns the underlying logrus logger
func GetLogrusLogger(l Logger) *logrus.Logger {
	if logrusLog, ok := l.(*logrusLogger); ok {
		return logrusLog.logger
	}
	// If it's not a logrusLogger, create a new one
	return logrus.New()
}