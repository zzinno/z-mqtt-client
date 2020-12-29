// @Title  utils.go
// @Description  工具函数
// @Author   loveward  2020/12/28 12:37
package client

import (
	"github.com/vmihailenco/msgpack"
)

func (z *ZClient) pub(topic string, msg ZMsg) {
	data, _ := msgpack.Marshal(msg)
	token := z.client.Publish(topic, 2, false, data)
	token.Wait()
	if token.Error() != nil {
		z.Logger.Error("Publish error:", token.Error())
	}
}

func (z *ZClient) sub(topic string) {
	token := z.client.Subscribe(topic, 1, nil)
	token.Wait()
	if token.Error() != nil {
		z.Logger.Error("Subscribe error:", token.Error())
	} else {
		z.Logger.Info("Subscribed to topic:", topic)
	}
}

func (z *ZClient) checkErr(err error) {
	if err != nil {
		z.Logger.Error(err)
	}
}
