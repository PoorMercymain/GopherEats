package logger

import "go.uber.org/zap"

func ProvideLogger() (*zap.SugaredLogger, error) {
	const logFileName = "logfile.log"
	const stdout = "stdout"

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFileName, stdout}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}