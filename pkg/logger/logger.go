package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info 记录信息日志
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

// Infof 格式化记录信息日志
func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Error 记录错误日志
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

// Errorf 格式化记录错误日志
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}
