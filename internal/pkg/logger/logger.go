// Package logger implements wrapper for zap logger.
package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_ ILogger = (*zap.SugaredLogger)(nil)
)

type ILogger interface {
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	DPanicln(args ...interface{})
	Panicln(args ...interface{})
	Fatalln(args ...interface{})
}

type logger struct {
	logger ILogger
	level  zapcore.Level
	file   string
	*sync.Once
}

var log = &logger{file: "logfile.log", Once: &sync.Once{}, level: zapcore.InfoLevel}

func Logger() ILogger { // actually, it may be better to use specific type, not an interface
	var err error
	log.Once.Do(func() {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(log.level)
		cfg.OutputPaths = []string{log.file, "stdout"}

		var logger *zap.Logger
		logger, err = cfg.Build()

		log.logger = logger.Sugar()
	})

	if err != nil {
		panicStr := fmt.Sprintf("cannot proceed: couldn't initialize logger because of an error - %v", err)
		panic(panicStr)
	}

	return log.logger
}
