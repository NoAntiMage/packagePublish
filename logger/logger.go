package logger

import (
	"PackageServer/config"
	"fmt"
	"path"
	"runtime"

	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var FileAndStdoutWriter io.Writer

// 初始化logger日志记录器
func InitLog() {
	writers := []io.Writer{
		os.Stdout,
	}

	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})

	if config.ServerConf.LoggerMode != "" {
		logLevel, _ := logrus.ParseLevel(config.ServerConf.LoggerMode)
		Log.SetLevel(logLevel)
	} else {
		Log.SetLevel(logrus.DebugLevel)
	}

	LogDir := fmt.Sprintf("%v/%v", config.ServerConf.DataDir, config.ServerConf.Name)
	if err := os.MkdirAll(LogDir, os.ModePerm); err != nil {
		fmt.Printf("logDir create fail %v", err)
	}
	logFileLocation := fmt.Sprintf("%v/access.log", LogDir)
	file, err := os.OpenFile(logFileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to log to file, using default stderr.")
	} else {
		writers = append(writers, file)
	}
	FileAndStdoutWriter = io.MultiWriter(writers...)
	Log.SetOutput(FileAndStdoutWriter)
}
