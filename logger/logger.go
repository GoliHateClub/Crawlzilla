package logger

import (
	"Crawlzilla/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type WriteSyncerMap struct {
	sync.Mutex
	syncers map[string]zapcore.WriteSyncer
}

func NewWriteSyncerMap() *WriteSyncerMap {
	return &WriteSyncerMap{
		syncers: make(map[string]zapcore.WriteSyncer),
	}
}

func (wsm *WriteSyncerMap) GetWriteSyncer(scope string) zapcore.WriteSyncer {
	wsm.Lock()
	defer wsm.Unlock()

	if ws, ok := wsm.syncers[scope]; ok {
		return ws
	}

	_, b, _, _ := runtime.Caller(0)

	root := filepath.Join(filepath.Dir(b), "..")

	file, err := os.OpenFile(fmt.Sprintf("%s/logs/%s.log", root, scope), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	ws := zapcore.AddSync(file)
	wsm.syncers[scope] = ws
	return ws
}

type ConfigLoggerType = func(scope string) (*zap.Logger, error)

func CreateLogger(scopes ...string) ConfigLoggerType {
	encoderCfg := zap.NewProductionEncoderConfig()
	isDev := config.GetBoolean("DEV_MODE")

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	writeSyncerMap := NewWriteSyncerMap()

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       isDev,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       nil,
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	// Create a base logger
	baseLogger := zap.Must(zapConfig.Build())

	// Create a composite logger core
	loggers := make(map[string]*zap.Logger)

	for _, scope := range scopes {
		loggers[scope] = baseLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(
				c,
				zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderCfg),
					writeSyncerMap.GetWriteSyncer(scope),
					zap.InfoLevel,
				),
			)
		}))
	}

	return func(scope string) (*zap.Logger, error) {
		_, ok := loggers[scope]

		if !ok {
			return nil, errors.New(fmt.Sprintf("%s scope is not exist", scope))
		}

		return loggers[scope].With(zap.String("scope", scope)), nil
	}
}

func ConfigLogger() ConfigLoggerType {
	logConfig := CreateLogger("crawls", "bot", "database")
	return logConfig
}
