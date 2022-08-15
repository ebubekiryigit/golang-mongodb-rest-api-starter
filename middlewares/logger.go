package middlewares

import (
	"io"
	"os"
	"path"
)

const LogPath = "logs/"
const LogFile = "access.log"

func LogWriter() io.Writer {
	_ = os.Mkdir(LogPath, 0770)

	logFilePath := path.Join(LogPath, LogFile)

	logFile, _ := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	return io.MultiWriter(logFile, os.Stdout)
}
