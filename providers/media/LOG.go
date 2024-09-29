package media

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

// StartLogger starts logger and sets it to be ready to used
func StartLogger() *Logger {

	logger := logrus.New()

	// Set the output of the logger to stdout
	logger.SetOutput(os.Stdout)

	// Set the formatter to TextFormatter with colored output
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)
	logger.AddHook(&callerHook{})

	return &Logger{logger: logger}
}

// InfoLogger logs message to the terminal
func (l *Logger) InfoLogger(message ...string) {
	l.logger.Info(fmt.Sprintf("::: ----> %v", message))
}

// WarningLogger logs a warning log to the teminal and saves it to a .log file
func (l *Logger) WarningLogger(message ...string) {
	l.logger.Warn(message)
}

// ErrorLog logs an error to the terminal log and saves it to a .log file
func (l *Logger) ErrorLog(message ...string) {
	l.logger.Error(message)
}

// FatalLog logas an Fatal log in terminal, saves the log into a .log file and exits program with exit status 1
func (l *Logger) FatalLog(message ...string) {
	defer l.logger.Fatal(message)
}

// callerHook adjusts the caller information to point to the correct file and line
type callerHook struct{}

func (hook *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *callerHook) Fire(entry *logrus.Entry) error {
	if entry.Caller != nil {
		frame := getCallerFrame(10)
		entry.Caller = &frame
	}
	return nil
}

// getCallerFrame retrieves the runtime frame for the given skip
func getCallerFrame(skip int) runtime.Frame {
	pc := make([]uintptr, 10)
	runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return frame
}
