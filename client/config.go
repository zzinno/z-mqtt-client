// @Title  config.go
// @Description  配置结构体
// @Author   loveward  2020/12/28 12:36
package client

import "github.com/zzinno/z-mqtt-client/logger"

//
type Config struct {
	Broker    string
	Port      int
	ClientID  string
	Username  string
	Password  string
	SelfTopic string
	CallBack  func(msg *RequestMsg) (*[]byte, error)
	Logger    logger.Logger
}
