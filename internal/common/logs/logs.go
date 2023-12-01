package logs

import (
	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger
type Fields = logrus.Fields

const (
	DebugLevel = logrus.DebugLevel
)

func Init(isDebug bool) {
	if !isDebug {
		SetFormatter(logrus.StandardLogger())
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
}

func New(isDebug bool) *Logger {
	logger := logrus.New()

	if !isDebug {
		SetFormatter(logger)
	}

	return logger
}

func GetLogger() *Logger {
	return logrus.StandardLogger()
}

func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func Warningln(args ...interface{}) {
	logrus.Warningln(args...)
}

func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func Print(i ...interface{}) {
	logrus.Print(i...)
}

func Printf(s2 string, i ...interface{}) {
	logrus.Printf(s2, i...)
}

func Println(i ...interface{}) {
	logrus.Println(i...)
}

func Fatal(i ...interface{}) {
	logrus.Fatal(i...)
}

func Fatalf(s2 string, i ...interface{}) {
	logrus.Fatalf(s2, i...)
}

func Fatalln(i ...interface{}) {
	logrus.Fatalln(i...)
}

func Panic(i ...interface{}) {
	logrus.Panic(i...)
}

func Panicf(s2 string, i ...interface{}) {
	logrus.Panicf(s2, i...)
}

func Panicln(i ...interface{}) {
	logrus.Panicln(i...)
}
