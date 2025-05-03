package logger

import (
	"go.uber.org/zap"
	"log"
)

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer Log.Sync()
}


// Logger defines the interface
type Logger interface {
	Info(msg string)
	Error(msg string)
}

// StdLogger implements Logger using the standard library
type StdLogger struct{}

func (l *StdLogger) Info(msg string) {
	log.Println("[INFO]", msg)
}

func (l *StdLogger) Error(msg string) {
	log.Println("[ERROR]", msg)
}

// New returns a new instance of StdLogger
func New() Logger {
	return &StdLogger{}
}
