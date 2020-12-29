// @Title  client.go
// @Description  客户端
// @Author   loveward  2020/12/28 12:33
package client

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack"
	"github.com/zzinno/z-mqtt-client/logger"
	"runtime"
	"sync"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type ZClient struct {
	messageHandler     mqtt.MessageHandler
	connectHandler     mqtt.OnConnectHandler
	connectLostHandler mqtt.ConnectionLostHandler
	client             mqtt.Client
	Logger             logger.Logger
	topic              string
	CallBack           func(msg *RequestMsg) (*[]byte, error)
	msgMap             sync.Map
}

// 使用时只要传入config就行
func (z *ZClient) New(c Config) error {
	if z.Logger == nil {
		z.Logger = new(logger.ZMqttLogger)
	}
	z.CallBack = c.CallBack
	z.connectHandler = z.onConnect
	z.connectLostHandler = z.onConnectLost
	z.messageHandler = z.deal
	z.topic = c.SelfTopic
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.Broker, c.Port))
	opts.SetClientID(c.ClientID)
	opts.SetUsername(c.Username)
	opts.SetPassword(c.Password)
	opts.SetDefaultPublishHandler(z.messageHandler)
	opts.OnConnect = z.connectHandler
	opts.OnConnectionLost = z.connectLostHandler
	opts.AutoReconnect = true
	opts.WillQos = 2
	z.client = mqtt.NewClient(opts)
	if token := z.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	z.sub(c.SelfTopic)
	return nil
}
func (z *ZClient) onConnect(_ mqtt.Client) {
	z.Logger.Info("Connected")
}

func (z *ZClient) onConnectLost(_ mqtt.Client, err error) {
	z.Logger.Error("Connect lost:", err)
}

func (z *ZClient) deal(_ mqtt.Client, msg mqtt.Message) {
	// 此处处理收到的所有消息
	buf := new(ZMsg)
	err := msgpack.Unmarshal(msg.Payload(), buf)
	z.checkErr(err)
	switch buf.MsgType {
	case 0:
		z.dealRequest(&buf.MsgContent)

	case 1:
		z.dealResponse(&buf.MsgContent)
	default:
		return
	}
}

func (z *ZClient) dealRequest(request *[]byte) {
	req := new(RequestMsg)
	err := msgpack.Unmarshal(*request, req)
	z.checkErr(err)
	go z.Logger.Info("Received Request:", string(req.Data.Key), "from topic:", req.FromTopic)
	res := new(RespondMsg)
	res.RequestMsg = *req
	// 如果没有读取到
	if req.MsgId == "" {
		res.Err = errors.New("no MsgId Found")
	} else {
		data, err2 := z.CallBack(req)
		if err2 != nil {
			res.Err = err

		} else {
			res.Data = *data
		}
	}
	content, _ := msgpack.Marshal(res)
	msg := ZMsg{
		MsgType:    1,
		MsgContent: content,
	}
	z.pub(req.FromTopic, msg)

}

func (z *ZClient) dealResponse(response *[]byte) {
	res := new(RespondMsg)
	err := msgpack.Unmarshal(*response, res)
	z.checkErr(err)
	go z.Logger.Info("Received Response:", string(res.Data))
	c, success := z.msgMap.LoadAndDelete(res.RequestMsg.MsgId)
	if success != true {
		z.Logger.Error("SyncMap Action Error")
	} else {
		ch := c.(*chan RespondMsg)
		*ch <- *res
	}

}

func (z *ZClient) SendMsgAndWaitReply(data Data, topic string, waitTime int) (RespondMsg, error) {
	id := uuid.New().String()
	req := RequestMsg{
		MsgId:     id,
		Data:      data,
		FromTopic: z.topic,
	}
	Content, _ := msgpack.Marshal(req)
	msg := ZMsg{
		MsgType:    0,
		MsgContent: Content,
	}
	z.pub(topic, msg)
	ch := make(chan RespondMsg, 1)
	timeOut := make(chan bool, 1)
	z.msgMap.Store(id, &ch)
	complete := false
	defer func() {
		z.msgMap.Delete(id)
		close(ch)
		close(timeOut)
	}()

	go func(c *bool, tc *chan bool) {
		time.Sleep(time.Duration(waitTime) * time.Second) // 等待n秒钟
		if complete == false {
			timeOut <- true
		}
	}(&complete, &timeOut)
	select {
	case m := <-ch:
		{
			complete = true
			return m, nil
		}
	case <-timeOut:
		{
			go z.Logger.Error("client waiting reply time out:", req)
			return RespondMsg{}, errors.New("client waiting reply time out")
		}
	}

}

func (z *ZClient) Loop() {
	select {}
}
