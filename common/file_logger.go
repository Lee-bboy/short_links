package common

import (
	"github.com/aiwuTech/fileLogger"
)

var (
	logFile *fileLogger.FileLogger
	INFO    *fileLogger.FileLogger
	WARN    *fileLogger.FileLogger
	ERROR   *fileLogger.FileLogger
)

func init() {
	INFO = fileLogger.NewSizeLogger("log", "info.log", "[WARN] ", 1000, 500, fileLogger.DEFAULT_FILE_UNIT, fileLogger.DEFAULT_LOG_SCAN, fileLogger.DEFAULT_LOG_SEQ)
	WARN = fileLogger.NewSizeLogger("log", "warn.log", "[WARN] ", 1000, 500, fileLogger.DEFAULT_FILE_UNIT, fileLogger.DEFAULT_LOG_SCAN, fileLogger.DEFAULT_LOG_SEQ)
	ERROR = fileLogger.NewSizeLogger("log", "error.log", "[ERROR] ", 1000, 500, fileLogger.DEFAULT_FILE_UNIT, fileLogger.DEFAULT_LOG_SCAN, fileLogger.DEFAULT_LOG_SEQ)
}
