package logger

import (
	"encoding/json"
	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	rawJSONConfig := []byte(`{
      "level": "info",
      "encoding": "console",
      "outputPaths": ["stdout"],
      "encoderConfig": {
        "messageKey": "message",
        "levelKey": "level",
        "nameKey": "logger",
        "timeKey": "time",
        "stacktraceKey": "stacktrace",
        "callstackKey": "callstack",
        "errorKey": "error",
        "timeEncoder": "iso8601",
        "fileKey": "file",
        "levelEncoder": "capitalColor",
        "durationEncoder": "second",
        "callerEncoder": "full",
        "nameEncoder": "full",
        "sampling": {
            "initial": "3",
            "thereafter": "10"
        }
      }
    }`)

	config := zap.Config{}
	if err := json.Unmarshal(rawJSONConfig, &config); err != nil {
		panic(err)
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}
