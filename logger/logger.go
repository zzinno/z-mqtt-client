// @Title  logger.go
// @Description
// @Author   loveward  2020/12/28 13:08
package logger

import (
	"log"
	"time"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
}
type ZMqttLogger struct{}

func (l ZMqttLogger) Info(args ...interface{}) {
	l.log("INFO", args)
}
func (l ZMqttLogger) Error(args ...interface{}) {
	l.log("ERROR", args)
}
func (l ZMqttLogger) Warning(args ...interface{}) {
	l.log("WARN", args)
}
func (l ZMqttLogger) log(T string, args ...interface{}) {
	currentTime := time.Now()
	log.Println(currentTime, "["+T+"]", args)
}
